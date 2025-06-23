package dao

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"srv/inventorySrv/model"
)

type InventoryDao struct {
	db *gorm.DB
}

type Stock struct {
	GoodsId int32
	Stocks  int32
}

func NewInventoryDao(db *gorm.DB) *InventoryDao {
	return &InventoryDao{
		db: db,
	}
}

func (d *InventoryDao) StockDecrease(decr *[]Stock) error {
	// update mm_sku set stock = case when id = 1 then stock - 1 when id = 2 then stock - 2 end where id in (1,2)
	if len(*decr) == 0 {
		return errors.New("商品库存扣减失败")
	}
	str := " WHEN goods_id = %d THEN stocks - %d "
	when := ""
	orStr := "(goods_id = %d AND stocks >= %d)"
	where := ""
	for index, item := range *decr {
		when += fmt.Sprintf(str, item.GoodsId, item.Stocks)
		where += fmt.Sprintf(orStr, item.GoodsId, item.Stocks)
		if index < len(*decr)-1 {
			where += "OR"
		}
	}

	sql := `UPDATE inventory SET stocks = CASE %s END WHERE %s`
	sql = fmt.Sprintf(sql, when, where)
	//res := &model.MmSku{}
	res := d.db.Model(&model.Inventory{}).Exec(sql)

	if res.RowsAffected != int64(len(*decr)) {
		return errors.New("商品库存不足")
	}

	return nil
}

func (d *InventoryDao) StockIncrease(decr *[]Stock) error {
	if len(*decr) == 0 {
		return nil
	}
	// 没有问题，这是不像是扣减需要防呆
	str := " WHEN goods_id = %d THEN stocks + %d "
	when := ""
	idx := ""
	for index, update := range *decr {
		when += fmt.Sprintf(str, update.GoodsId, update.Stocks)
		idx += fmt.Sprintf("%d", update.GoodsId)
		if index < len(*decr)-1 {
			idx += ","
		}
	}
	sql := `UPDATE inventory SET stocks = CASE %s END WHERE id IN (%s) `
	sql = fmt.Sprintf(sql, when, idx)
	tx := d.db.Model(&model.Inventory{}).Exec(sql)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
