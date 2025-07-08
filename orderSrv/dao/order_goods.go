package dao

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	"srv/orderSrv/model"
)

type OrderGoodsDao struct {
	db *gorm.DB
}

func NewOrderGoodsDao(db *gorm.DB) *OrderGoodsDao {
	return &OrderGoodsDao{db: db}
}

func (r *OrderGoodsDao) BatchCreate(ctx context.Context, param []model.OrderGoods) error {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("dao", "BatchCreate"))
	defer span.End()

	res := r.db.Model(&model.OrderGoods{}).CreateInBatches(param, len(param))
	return HandleError(res, len(param))
}
