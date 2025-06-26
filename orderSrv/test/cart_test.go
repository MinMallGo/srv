package test

import (
	"context"
	"log"
	proto "srv/orderSrv/proto/gen"
	"testing"
)

func TestCartCrud(t *testing.T) {
	c := proto.NewCartClient(SrvInit())
	defer SrvConn.Close()
	_, err := c.AddGoods(context.Background(), &proto.AddGoodsReq{
		GoodsId:  2,
		GoodsNum: 2,
		GoodsImg: "http://example.com/goods/phone_front01.jpg",
		UserId:   2,
	})
	if err != nil {
		log.Panicln("RemoveGoods: GoodsId:  2,\n\t\tGoodsNum: 2,", err.Error())
	}

	if _, err = c.AddGoods(context.Background(), &proto.AddGoodsReq{
		GoodsId:  1,
		GoodsNum: 2,
		GoodsImg: "http://example.com/goods/phone_front02.jpg",
		UserId:   2,
	}); err != nil {
		log.Panicln("AddGoods: AddGoods(context.Background(), &proto.AddGoodsReq{\n\t\tGoodsId:  1,\n\t\tGoodsNum: 2,\n\t\tGoodsImg: \"http://example.com/goods/phone_front02.jpg\",\n\t\tUserId:   2,", err.Error())
	}

	if _, err = c.AddGoods(context.Background(), &proto.AddGoodsReq{
		GoodsId:  1,
		GoodsNum: 2,
		GoodsImg: "http://example.com/goods/phone_front02.jpg",
		UserId:   2,
	}); err != nil {
		log.Panicln("RemoveGoods: GoodsId:  1,\n\t\tGoodsNum: 2,", err.Error())
	}

	if _, err = c.AddGoods(context.Background(), &proto.AddGoodsReq{
		GoodsId:  3,
		GoodsNum: 2,
		GoodsImg: "http://example.com/goods/phone_front02.jpg",
		UserId:   2,
	}); err != nil {
		log.Panicln("RemoveGoods: GoodsId:  3,\n\t\tGoodsNum: 2,", err.Error())
	}

	if _, err = c.RemoveGoods(context.Background(), &proto.RemoveGoodsReq{
		GoodsId: []int32{3},
		UserId:  2,
	}); err != nil {
		log.Panicln("RemoveGoods: ", err.Error())
	}

	if _, err = c.SelectGoods(context.Background(), &proto.SelectGoodsReq{
		GoodsId: []int32{1, 2},
		UserId:  2,
	}); err != nil {
		log.Panicln("SelectGoods: ", err.Error())
	}

	res, err := c.GetCartList(context.Background(), &proto.GetCartListReq{
		UserId: 2,
	})

	if res.Total != 2 {
		log.Panicln("GetCartList: ", err.Error())
	}

	if _, err = c.UpdateGoodsNum(context.Background(), &proto.UpdateNumReq{
		GoodsId:  2,
		GoodsNum: 10,
		UserId:   2,
	}); err != nil {
		log.Panicln("更新失败", err.Error())
	}

}
