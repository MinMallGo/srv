package test

import (
	"context"
	"log"
	proto "srv/orderSrv/proto/gen"
	"testing"
)

func TestOrderCreate(t *testing.T) {
	// 创建商品。先获取连接
	c := proto.NewOrderClient(SrvInit())
	defer SrvConn.Close()
	order, err := c.CreateOrder(context.Background(), &proto.CreateOrderReq{
		UserID:          2,
		Address:         "湖北长沙",
		RecipientName:   "小白",
		RecipientMobile: "17623456789",
		Message:         "尽快发货",
	})
	if err != nil {
		panic(err)
	}

	log.Println(order)
}

func TestOrderList(t *testing.T) {
	c := proto.NewOrderClient(SrvInit())
	defer SrvConn.Close()
	list, err := c.GetList(context.Background(), &proto.OrderListReq{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		panic(err)
	}
	log.Println(list)
}

func TestOrderDetail(t *testing.T) {
	c := proto.NewOrderClient(SrvInit())
	defer SrvConn.Close()
	list, err := c.GetListDetail(context.Background(), &proto.OrderDetailReq{
		OrderId: 3,
		OrderSn: "",
	})
	if err != nil {
		panic(err)
	}
	log.Println(list)
}
