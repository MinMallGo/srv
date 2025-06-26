package dao

import (
	"gorm.io/gorm"
	"srv/orderSrv/model"
)

type OrderGoodsDao struct {
	db *gorm.DB
}

func NewOrderGoodsDao(db *gorm.DB) *OrderGoodsDao {
	return &OrderGoodsDao{db: db}
}

func (r *OrderGoodsDao) BatchCreate(param []model.OrderGoods) error {
	res := r.db.Model(&model.OrderGoods{}).CreateInBatches(param, len(param))
	return HandleError(res, len(param))
}
