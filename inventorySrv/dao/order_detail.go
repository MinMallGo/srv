package dao

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"srv/inventorySrv/global"
	"srv/inventorySrv/model"
	"srv/inventorySrv/structs"
)

type OrderDetailDao struct {
	db *gorm.DB
}

func NewOrderDetailDao(db *gorm.DB) *OrderDetailDao {
	return &OrderDetailDao{
		db: global.DB,
	}
}

func (r *OrderDetailDao) Create(orderSN string, stocks []structs.Stocks) error {
	detail := make([]model.Details, 0, len(stocks))
	for _, stock := range stocks {
		detail = append(detail, model.Details{
			GoodsID:  stock.GoodsID,
			GoodsNum: stock.Stock,
		})
	}
	res := r.db.Model(&model.OrderHistory{}).Create(&model.OrderHistory{
		OrderSN: orderSN,
		Status:  1,
		Details: detail,
	})
	if res.Error != nil {
		zap.L().Error("创建订单详细失败", zap.Error(res.Error))
		return errors.New("创建订单详细失败")
	}

	if res.RowsAffected == 0 {
		return errors.New("创建订单详细失败")
	}
	return nil
}

func (r *OrderDetailDao) GetOne(orderSN string) (*model.OrderHistory, error) {
	resp := &model.OrderHistory{}
	x := r.db.Model(&model.OrderHistory{}).Where("order_sn = ?", orderSN).First(resp)
	if x.Error != nil {
		return resp, x.Error
	}

	if x.RowsAffected == 0 {
		return resp, gorm.ErrRecordNotFound
	}

	if resp.Status == 2 {
		return resp, errors.New("订单不存在")
	}

	return resp, nil
}

func (r *OrderDetailDao) UpdateStatus(orderSN string) error {
	resp := r.db.Model(&model.OrderHistory{}).Where("order_sn = ?", orderSN).Update("status", 2)
	if resp.Error != nil || resp.RowsAffected == 0 {
		return errors.New("更新归还状态失败")
	}
	
	return nil
}
