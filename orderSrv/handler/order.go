package handler

import (
	"context"
	dao2 "srv/orderSrv/dao"
	"srv/orderSrv/global"
	proto "srv/orderSrv/proto/gen"
)

type OrderService struct {
	proto.UnimplementedOrderServer
}

func (o OrderService) CreateOrder(ctx context.Context, req *proto.CreateOrderReq) (*proto.CreateResp, error) {
	//TODO implement me
	panic("implement me")
}

func (o OrderService) GetList(ctx context.Context, req *proto.OrderListReq) (*proto.OrderListResp, error) {
	res, err := dao2.NewOrderDao(global.DB).GetList(ctx, dao2.OrderListResp{Page: req.Page, Size: req.PageSize})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (o OrderService) GetListDetail(ctx context.Context, req *proto.OrderDetailReq) (*proto.OrderDetailResp, error) {
	//TODO implement me
	panic("implement me")
}

func (o OrderService) mustEmbedUnimplementedOrderServer() {
	//TODO implement me
	panic("implement me")
}
