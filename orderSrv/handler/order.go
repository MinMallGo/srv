package handler

import (
	"context"
	"fmt"
	rmq_client "github.com/apache/rocketmq-clients/golang/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	dao2 "srv/orderSrv/dao"
	"srv/orderSrv/global"
	"srv/orderSrv/model"
	proto "srv/orderSrv/proto/gen"
	"srv/orderSrv/service"
	"srv/orderSrv/utils"
)

type OrderMsg struct {
	OrderSN string      `json:"orderSN"`
	OrderID string      `json:"orderID"`
	UserID  string      `json:"userId"`
	Goods   []GoodsInfo `json:"goods"`
}

type GoodsInfo struct {
	GoodsID uint64 `json:"goodsID"`
	GoodsSN string `json:"goodsSN"`
}

type OrderService struct {
	proto.UnimplementedOrderServer
}

// TODO 这里要修改一下，变成不在这里进行库存的扣减，而是通过成功发送事务消息，在库存侧监听来达到目的
// 同时这样也方面扩展，有新的服务需要加入进来就订阅然后执行自己的任务就行

func (o OrderService) CreateOrder(ctx context.Context, req *proto.CreateOrderReq) (*proto.CreateResp, error) {
	tracer := otel.Tracer("test-tracer") // 配置
	var span trace.Span
	ctx, span = tracer.Start(ctx, fmt.Sprintf("CreateOrder"))
	defer span.End()

	// 跨服务调用，这里会用到库存服务，以及商品服务，
	// 通过引入其他服务的proto，然后连接客户端，然后再进行通信
	// TODO 分布式事务怎么保证一致性
	// 查询购物车里面勾选的物品
	// 	确实jd是读取购物车里面的勾选信息

	// 分布式事务怎么保证一致性 ： 消息队列的事务消息
	res, err := NewTrxMsg()
	if err != nil {
		return nil, err
	}
	producer := res.Product
	// TODO 还需要看下这里会有什么坑
	defer producer.GracefulStop()
	// 分布式事务怎么保证一致性 ： 消息队列的事务消息
	// 先生成订单唯一编号
	orderSN := global.OrderSN(int(req.UserID))
	zap.L().Info("<<orderSN>>", zap.String("orderSN", orderSN))
	// 这里是先发送一个办消息，然后后面的本地事务再发送 trx.Commit or trx.Rollback

	ctx, childSpan := tracer.Start(ctx, "SendHalf-Message")

	trx := producer.BeginTransaction()
	msg := &rmq_client.Message{
		Topic: TrxTopic,
		Body:  []byte(orderSN),
	}
	msg.SetTag("returnStock")
	msg.SetKeys("order_sn", orderSN)

	receipts, err := producer.SendWithTransaction(ctx, msg, trx)
	if err != nil {
		zap.L().Error("发送半消息失败", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "内部错误")
	}

	zap.L().Info("发送半消息成功", zap.Any("receipts", receipts))
	childSpan.End()

	ctx, cartSpan := tracer.Start(ctx, "GetCartList")
	list, err := dao2.NewCartDao(global.DB).CartList(ctx, dao2.CartBase{UserId: req.UserID})
	cartSpan.End()
	if err != nil {
		return nil, err
	}

	if len(*list) == 0 {
		return nil, status.Error(codes.NotFound, "请选择商品后购买")
	}

	goodsIds := make([]int32, 0, len(*list))
	goodsMap := make(map[int32]int32)
	sell := make([]*proto.SetInfo, 0, len(*list))
	for _, v := range *list {
		goodsIds = append(goodsIds, v.GoodsID)
		goodsMap[v.GoodsID] = v.Nums
		sell = append(sell, &proto.SetInfo{
			GoodsId: v.GoodsID,
			Stock:   v.Nums,
		})
	}

	ctx, batchSpan := tracer.Start(ctx, "GetGoodsBatch")
	// 查询商品服务，获取查看库存是否足够
	batchGoods, err := global.CrossSrv.Goods.BatchGetGoods(context.Background(), &proto.BatchGoodsInfo{
		Id: goodsIds,
	})
	if err != nil {
		zap.L().Error("查询商品服务失败", zap.Error(err))
		trx.RollBack()
		return nil, err
	}
	batchSpan.End()

	if len(batchGoods.Data) == 0 || len(batchGoods.Data) != len(goodsMap) {
		return nil, status.Error(codes.NotFound, "请选择正确的商品")
	}

	ctx, sellSpan := tracer.Start(ctx, "SellStock")
	// 库存扣减
	_, err = global.CrossSrv.Inventory.SellStock(context.Background(), &proto.MultipleInfo{OrderSN: orderSN, Sell: sell})
	sellSpan.End()
	if err != nil {
		zap.L().Error("库存扣减失败", zap.Error(err))
		trx.RollBack()
		return nil, err
	}

	// 再往下就是写入到order以及order_goods

	orderGoods := make([]model.OrderGoods, 0)
	var total_price float32 = 0
	subject := ""

	for _, datum := range batchGoods.Data {
		price := datum.ShopPrice * float32(goodsMap[datum.Id])
		total_price += price
		subject += ` ` + datum.Name + `; `
		img := ""
		if len(datum.ImageUrl) > 0 {
			img = datum.ImageUrl[0]
		}
		orderGoods = append(orderGoods, model.OrderGoods{
			OrderSN:    orderSN,
			GoodsId:    datum.Id,
			GoodsPrice: price, // 单价 * 数量
			PayPrice:   price, // TODO 优惠券之类的
			GoodsName:  datum.Name,
			Nums:       goodsMap[datum.Id],
			GoodsImg:   img,
		})
	}

	order := &model.Order{
		UserID:          req.UserID,
		OrderSN:         orderSN,
		Status:          "PAYING",
		SubjectTitle:    subject,
		OrderPrice:      total_price,
		FinalPrice:      total_price,
		Address:         req.Address,
		RecipientName:   req.RecipientName,
		RecipientMobile: req.RecipientMobile,
		Message:         req.Message,
	}

	var orderId int
	err = global.DB.Transaction(func(tx *gorm.DB) error {
		orderId, err = dao2.NewOrderDao(tx).Create(ctx, order)
		zap.L().Info(fmt.Sprintf("Create order ID is :%d", orderId))
		if err != nil {
			return err
		}
		for i := range orderGoods {
			orderGoods[i].OrderId = int32(orderId)
		}
		// 批量创建订单商品
		if err = dao2.NewOrderGoodsDao(tx).BatchCreate(ctx, orderGoods); err != nil {
			return err
		}

		ctx, cartSpan = tracer.Start(ctx, "deleteCart")
		// 从购物车里面移除购买的商品
		if err = dao2.NewCartDao(tx).Delete(context.Background(), dao2.CartMultiGoods{
			UserId:  req.UserID,
			GoodsId: goodsIds,
		}); err != nil {
			return err
		}
		cartSpan.End()

		// 这里加一个发送订单取消的延时消息
		if err = SendOrderDelayMsg(orderSN); err != nil {
			return err
		}

		return nil
	})
	// 这里，判断事务是否成功，err != nil 就需要使用
	if err != nil {
		zap.L().Error("本地事务提交失败", zap.Error(err))
		trx.RollBack()
		return nil, err
	}
	// 提交全消息
	trx.Commit()

	return &proto.CreateResp{
		OrderId: int32(orderId),
		OrderSn: orderSN,
	}, nil
}

// 正常的半消息事务应该是这样的：
// 本地事务在启动之前，就启动一个半事务消息，然后本地事务完成之后，提交这个半消息
// 然后订单创建的本地应该：订阅消息订单回滚的消息队列，然后收到消息之后进行回滚操作
// 商品库存侧：订阅这个半事务消息的队列，然后进行业务操作，同时需要监听一个库存回滚的队列，然后进行库存回滚
// 其他业务应该也是。
// todo 改一下吧，这里改成只调用，而不是写太多的业务逻辑在这里面

func (o OrderService) CreateOrder2(ctx context.Context, req *proto.CreateOrderReq) (*proto.CreateResp, error) {
	// 需要解耦一下代码
	// 将发送消息封装一下
	// 将链路追踪封装一下
	// 封装一下事务发送
	tracer := otel.Tracer("order-tracer") // 配置
	ctx, span := tracer.Start(ctx, "CreateOrder")
	defer span.End()

	return utils.WithSpan[*proto.CreateResp](ctx, tracer, "create_order", []attribute.KeyValue{
		attribute.Int("user_id", int(req.UserID)),
	}, func(ctx context.Context) (*proto.CreateResp, error) {
		// 查询消息调用service里面的内容
		return service.CreateOrder(ctx, tracer, req)
	})
}

func (o OrderService) GetList(ctx context.Context, req *proto.OrderListReq) (*proto.OrderListResp, error) {
	tracer := otel.Tracer("test-tracer") // 配置
	var span trace.Span
	ctx, span = tracer.Start(ctx, fmt.Sprintf("span-OrderService.GetList"))
	defer span.End()
	//span := trace.SpanFromContext(ctx)
	//span.SetAttributes(attribute.String("service.name", "<GetList>"))
	//defer span.End()

	res, err := dao2.NewOrderDao(global.DB).GetList(ctx, dao2.OrderListReq{Page: req.Page, Size: req.PageSize, UserID: req.UserId})
	//iSpan.End()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (o OrderService) GetListDetail(ctx context.Context, req *proto.OrderDetailReq) (*proto.OrderDetailResp, error) {
	tracer := otel.Tracer("test-tracer") // 配置
	var iSpan trace.Span
	ctx, iSpan = tracer.Start(ctx, fmt.Sprintf("span-OrderService.GetListDetail"))
	// 参数检查
	if req.OrderId == 0 && len(req.OrderSn) == 0 {
		return nil, status.Error(codes.InvalidArgument, "参数异常")
	}

	res, err := dao2.NewOrderDao(global.DB).GetDetail(ctx, dao2.OrderDetailReq{OrderId: req.OrderId, OrderSN: req.OrderSn})
	if err != nil {
		return nil, err
	}

	data := make([]*proto.GoodsInfo, 0, len(res.Goods))
	for _, goods := range res.Goods {
		data = append(data, &proto.GoodsInfo{
			OrderID:    goods.OrderId,
			OrderSN:    goods.OrderSN,
			GoodsID:    goods.GoodsId,
			GoodsPrice: goods.GoodsPrice,
			PayPrice:   goods.PayPrice,
			GoodsName:  goods.GoodsName,
			Num:        goods.Nums,
		})
	}

	resp := &proto.OrderDetailResp{
		UserID:          res.UserID,
		OrderSN:         res.OrderSN,
		PayType:         res.PayType,
		Status:          res.Status,
		TradeNo:         res.TradeNo,
		SubjectTitle:    res.SubjectTitle,
		OrderPrice:      res.OrderPrice,
		FinalPrice:      res.FinalPrice,
		Address:         res.Address,
		RecipientName:   res.RecipientName,
		RecipientMobile: res.RecipientMobile,
		Message:         res.Message,
		Snapshot:        res.Snapshot,
		Goods:           data,
	}
	iSpan.End()

	return resp, nil
}

func (o OrderService) mustEmbedUnimplementedOrderServer() {
	//TODO implement me
	panic("implement me")
}
