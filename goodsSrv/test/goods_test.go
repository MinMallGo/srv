package test

import (
	"context"
	"encoding/json"
	proto "goodsSrv/proto/gen"
	"log"
	"testing"
)

func TestGoodsList(t *testing.T) {
	s := proto.NewGoodsClient(SrvInit())
	defer SrvClose()
	lists, err := s.GoodsList(context.Background(), &proto.GoodsFilterRequest{
		PriceMin:    0,
		PriceMax:    0,
		IsHot:       false,
		IsNew:       false,
		IsTab:       false,
		TopCategory: 3,
		Pages:       0,
		PageSize:    0,
		KeyWord:     "",
		Brand:       5,
	})
	if err != nil {
		t.Error(err)
	}
	log.Println(lists.Total)
	log.Println("========================")
	for _, list := range lists.Data {
		log.Println(list.Name, list.ShopPrice)
	}
}

func TestGoodsDetail(t *testing.T) {
	s := proto.NewGoodsClient(SrvInit())
	defer SrvClose()
	detail, err := s.GetGoodsDetail(context.Background(), &proto.GoodsInfoRequest{
		Id: 1,
	})
	if err != nil {
		t.Error(err)
	}
	log.Println(detail)
}

func TestGoodsCRUD(t *testing.T) {
	s := proto.NewGoodsClient(SrvInit())
	defer SrvClose()
	goods, err := s.CreateGoods(context.Background(), &proto.CreateGoodsInfo{
		CategoryId:      1,
		BrandId:         1,
		OnSale:          false,
		ShipFree:        false,
		IsNew:           false,
		Stock:           1999,
		Name:            "xafs",
		GoodsSn:         "asdf",
		ClickNum:        1,
		SoldNum:         1,
		FavNum:          1,
		MarketPrice:     1,
		ShopPrice:       1,
		GoodsBrief:      "123",
		ImageUrl:        []string{"123", "sadf"},
		Description:     []string{"123", "sadf"},
		GoodsFrontImage: "asdfasdf",
	})
	if err != nil {
		log.Panicln("创建失败：", err.Error())
	}

	log.Println(goods)

	_, err = s.UpdateGoods(context.Background(), &proto.UpdateGoodsInfo{
		Id:              goods.Id,
		CategoryId:      2,
		BrandId:         2,
		OnSale:          false,
		ShipFree:        false,
		IsNew:           false,
		Stock:           1999,
		Name:            "update_xafs",
		GoodsSn:         "update_xafs",
		ClickNum:        1,
		SoldNum:         1,
		FavNum:          1,
		MarketPrice:     1,
		ShopPrice:       1,
		GoodsBrief:      "123",
		ImageUrl:        []string{"123", "sadf"},
		Description:     []string{"123", "sadf"},
		GoodsFrontImage: "asdfasdf",
	})
	if err != nil {
		log.Panicln("<修改失败>", err.Error())
	}

	_, err = s.DeleteGoods(context.Background(), &proto.DeleteGoodsInfo{
		Id: goods.Id,
	})
	if err != nil {
		log.Panicln("<删除失败>", err.Error())
	}

}

func TestFuckMarshal(t *testing.T) {
	s := `["http://example.com/goods/dress_front01.jpg"]`
	a := make([]string, 0)
	if err := json.Unmarshal([]byte(s), &a); err != nil {
		t.Error(err)
	}
	log.Println(a)
}
