package test

import (
	"context"
	"encoding/json"
	"log"
	proto "srv/goodsSrv/proto/gen"
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

func TestFuckMarshal(t *testing.T) {
	s := `["http://example.com/goods/dress_front01.jpg"]`
	a := make([]string, 0)
	if err := json.Unmarshal([]byte(s), &a); err != nil {
		t.Error(err)
	}
	log.Println(a)
}
