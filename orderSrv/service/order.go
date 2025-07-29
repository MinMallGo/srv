package service

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	dao2 "srv/orderSrv/dao"
	"srv/orderSrv/global"
	"srv/orderSrv/model"
	proto "srv/orderSrv/proto/gen"
	"srv/orderSrv/rpc"
	"srv/orderSrv/utils"
)

func CreateOrder(ctx context.Context, tracer trace.Tracer, req *proto.CreateOrderReq) (*proto.CreateResp, error) {
	// 获取购物车列表
	list, err := GetCartList(ctx, tracer, int(req.UserID))
	if err != nil {
		return nil, err
	}

	goodsIds := make([]int32, 0, len(list))
	for _, info := range list {
		goodsIds = append(goodsIds, info.GoodsId)
	}

	// rpc调用获取商品信息
	goods, err := rpc.GetGoods(ctx, tracer, goodsIds)
	if err != nil {
		return nil, err
	}
	// 发送半事务消息
	res, err := utils.NewTrxMsg()
	if err != nil {
		return nil, err
	}
	// 最后优雅关闭消息
	defer res.Product.GracefulStop()
	res, err = utils.SendTrxMsg(ctx, res.Product, "", "")
	if err != nil {
		return nil, err
	}

	// 获取订单信息
	orderSN := global.OrderSN(int(req.UserID))
	// 整理创建的订单
	orderData := TidyOrder(goods, req, list, orderSN)
	orderGoodsData := TidyOrderGoods(goods, req, list, orderSN)

	var orderId int
	err = global.DB.Transaction(func(tx *gorm.DB) error {
		orderId, err = dao2.NewOrderDao(tx).Create(ctx, &orderData)
		if err != nil {
			return err
		}

		err = dao2.NewOrderGoodsDao(tx).BatchCreate(ctx, orderGoodsData)
		if err != nil {
			return err
		}

		if err = dao2.NewCartDao(tx).Delete(context.Background(), dao2.CartMultiGoods{
			UserId:  req.UserID,
			GoodsId: goodsIds,
		}); err != nil {
			return err
		}

		return nil
	})

	// 整理订单创建信息，然后提交事务
	if err != nil {
		res.Trx.RollBack()
		return nil, status.Error(codes.Internal, err.Error())
	}

	res.Trx.Commit()

	return &proto.CreateResp{
		OrderId: int32(orderId),
		OrderSn: orderSN,
	}, nil
}

// GetCartList 获取购物车列表
func GetCartList(ctx context.Context, tracer trace.Tracer, userID int) ([]*proto.SetInfo, error) {
	return utils.WithSpan[[]*proto.SetInfo](ctx, tracer, "GetCartList", []attribute.KeyValue{
		attribute.String("user_id", ""),
	}, func(ctx context.Context) ([]*proto.SetInfo, error) {
		list, err := dao2.NewCartDao(global.DB).CartList(ctx, dao2.CartBase{UserId: int32(userID)})
		if err != nil {
			return nil, err
		}

		if len(*list) == 0 {
			return nil, status.Error(codes.NotFound, "请选择商品后购买")
		}

		sell := make([]*proto.SetInfo, 0, len(*list))
		for _, v := range *list {
			sell = append(sell, &proto.SetInfo{
				GoodsId: v.GoodsID,
				Stock:   v.Nums,
			})
		}

		return sell, nil

	})
}

// TidyOrder 整理创建订单的数据
func TidyOrder(goods *proto.GoodsListResponse, req *proto.CreateOrderReq, cart []*proto.SetInfo, orderSN string) model.Order {
	var totalPrice float32 = 0
	subject := ""
	goodsMap := make(map[int32]int32)
	for _, info := range cart {
		goodsMap[info.GoodsId] = info.GoodsId
	}

	for _, datum := range goods.Data {
		price := datum.ShopPrice * float32(goodsMap[datum.Id])
		totalPrice += price
		subject += ` ` + datum.Name + `; `
	}

	return model.Order{
		OrderSN:         orderSN,
		UserID:          req.UserID,
		Status:          "PAYING",
		SubjectTitle:    subject,
		OrderPrice:      totalPrice,
		FinalPrice:      totalPrice,
		Address:         req.Address,
		RecipientName:   req.RecipientName,
		RecipientMobile: req.RecipientMobile,
		Message:         req.Message,
	}
}

// TidyOrderGoods 整理创建订单商品的数据
func TidyOrderGoods(goods *proto.GoodsListResponse, req *proto.CreateOrderReq, cart []*proto.SetInfo, orderSN string) []model.OrderGoods {
	var totalPrice float32 = 0
	subject := ""
	goodsMap := make(map[int32]int32)
	orderGoods := make([]model.OrderGoods, 0)
	for _, info := range cart {
		goodsMap[info.GoodsId] = info.GoodsId
	}

	for _, datum := range goods.Data {
		price := datum.ShopPrice * float32(goodsMap[datum.Id])
		totalPrice += price
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
	return orderGoods
}
