package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"log"
	dao2 "srv/orderSrv/dao"
	"srv/orderSrv/global"
	"srv/orderSrv/model"
	proto "srv/orderSrv/proto/gen"
)

type OrderService struct {
	proto.UnimplementedOrderServer
}

func (o OrderService) CreateOrder(ctx context.Context, req *proto.CreateOrderReq) (*proto.CreateResp, error) {
	// 跨服务调用，这里会用到库存服务，以及商品服务，
	// 通过引入其他服务的proto，然后连接客户端，然后再进行通信
	// TODO 分布式事务怎么保证一致性
	// 查询购物车里面勾选的物品
	// 	确实jd是读取购物车里面的勾选信息
	list, err := dao2.NewCartDao(global.DB).CartList(ctx, dao2.CartBase{UserId: req.UserID})
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
	// 查询商品服务，获取查看库存是否足够
	batchGoods, err := global.CrossSrv.Goods.BatchGetGoods(context.Background(), &proto.BatchGoodsInfo{
		Id: goodsIds,
	})
	if err != nil {
		return nil, err
	}

	if len(batchGoods.Data) == 0 || len(batchGoods.Data) != len(goodsMap) {
		return nil, status.Error(codes.NotFound, "请选择正确的商品")
	}

	// 库存扣减
	_, err = global.CrossSrv.Inventory.SellStock(context.Background(), &proto.MultipleInfo{Sell: sell})
	if err != nil {
		return nil, err
	}

	// 再往下就是写入到order以及order_goods

	orderGoods := make([]model.OrderGoods, 0)
	var total_price float32 = 0
	subject := ""
	orderSN := global.OrderSN(int(req.UserID))
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
		orderId, err = dao2.NewOrderDao(tx).Create(order)
		log.Println("Create order ID is :", orderId)
		if err != nil {
			return err
		}
		for i := range orderGoods {
			orderGoods[i].OrderId = int32(orderId)
		}

		err = dao2.NewOrderGoodsDao(tx).BatchCreate(orderGoods)
		if err != nil {
			return err
		}

		// 从购物车里面移除购买的商品
		err = dao2.NewCartDao(tx).Delete(context.Background(), dao2.CartMultiGoods{
			UserId:  req.UserID,
			GoodsId: goodsIds,
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &proto.CreateResp{
		OrderId: int32(orderId),
		OrderSn: orderSN,
	}, nil
}

func (o OrderService) GetList(ctx context.Context, req *proto.OrderListReq) (*proto.OrderListResp, error) {
	res, err := dao2.NewOrderDao(global.DB).GetList(ctx, dao2.OrderListReq{Page: req.Page, Size: req.PageSize, UserID: req.UserId})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (o OrderService) GetListDetail(ctx context.Context, req *proto.OrderDetailReq) (*proto.OrderDetailResp, error) {
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

	return resp, nil
}

func (o OrderService) mustEmbedUnimplementedOrderServer() {
	//TODO implement me
	panic("implement me")
}
