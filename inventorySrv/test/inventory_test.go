package test

import (
	"context"
	"fmt"
	proto "srv/inventorySrv/proto/gen"
	"sync"
	"testing"
)

func TestSet(t *testing.T) {
	_, err := proto.NewInventoryClient(SrvInit()).SetStock(context.Background(), &proto.SetInfo{
		GoodsId: 2,
		Stock:   299,
	})
	if err != nil {
		t.Error(err)
	}
}

func TestGet(t *testing.T) {
	e, err := proto.NewInventoryClient(SrvInit()).GetStock(context.Background(), &proto.GetInfo{
		GoodsId: 1,
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(e)
}

func TestSellStock(t *testing.T) {
	fn := func() {
		_, err := proto.NewInventoryClient(SrvInit()).SellStock(context.Background(), &proto.MultipleInfo{
			Sell: []*proto.SetInfo{
				{GoodsId: 1, Stock: 1},
				{GoodsId: 2, Stock: 1},
			},
		})
		if err != nil {
			t.Error(err)
		}
	}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fn()
		}()
	}
	wg.Wait()

}

func TestReturnStock(t *testing.T) {
	e, err := proto.NewInventoryClient(SrvInit()).ReturnStock(context.Background(), &proto.MultipleInfo{
		Sell: []*proto.SetInfo{
			{GoodsId: 1, Stock: 2},
			{GoodsId: 2, Stock: 2},
		},
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(e)
}
