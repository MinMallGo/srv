package rpc

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"srv/orderSrv/global"
	proto "srv/orderSrv/proto/gen"
	"srv/orderSrv/utils"
	"strconv"
	"strings"
)

func GetGoods(ctx context.Context, trace trace.Tracer, goodsIds []int32) (*proto.GoodsListResponse, error) {
	var s strings.Builder
	for i, id := range goodsIds {
		s.WriteString(strconv.Itoa(int(id)))
		if i != len(goodsIds)-1 {
			s.WriteString(",")
		}
	}

	return utils.WithSpan[*proto.GoodsListResponse](ctx, trace, "rpc_call_batch-goods", []attribute.KeyValue{
		attribute.String("goods_ids", s.String()),
	}, func(ctx context.Context) (*proto.GoodsListResponse, error) {
		batchGoods, err := global.CrossSrv.Goods.BatchGetGoods(context.Background(), &proto.BatchGoodsInfo{
			Id: goodsIds,
		})
		return batchGoods, err
	})

}
