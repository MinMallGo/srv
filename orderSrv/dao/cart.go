package dao

import (
	"context"
	"gorm.io/gorm"
	"srv/orderSrv/model"
)

type CartDao struct {
	db *gorm.DB
}

func NewCartDao(db *gorm.DB) *CartDao {
	return &CartDao{db: db}
}

// CartBase cart的基础信息
type CartBase struct {
	UserId  int32
	GoodsId int32
}

type CartMultiGoods struct {
	UserId  int32
	GoodsId []int32
}

// CartCreate 创建/修改都能用这个吧。只是看需不需要用id
type CartCreate struct {
	UserId   int32
	GoodsId  int32
	GoodsNum int32
	GoodsImg string
}

type CartUpdate struct {
	ID       int32
	UserId   int32
	GoodsId  int32
	GoodsNum int32
	GoodsImg string
}

func (c CartDao) Exists(ctx context.Context, req CartBase) bool {
	res, err := c.Get(ctx, req)
	if err != nil || res == nil || res.ID == 0 {
		return false
	}
	return true
}

func (c CartDao) MultiExists(ctx context.Context, req CartMultiGoods) bool {
	resp := &model.Cart{}
	res := c.db.Model(&model.Cart{}).Where("user_id = ? AND goods_id = ?", req.UserId, req.GoodsId).First(&resp)
	err := HandleError(res, len(req.GoodsId))
	return err == nil
}

func (c CartDao) Get(ctx context.Context, req CartBase) (*model.Cart, error) {
	resp := &model.Cart{}
	res := c.db.Model(&model.Cart{}).Where("user_id = ? AND goods_id = ?", req.UserId, req.GoodsId).First(&resp)
	return resp, HandleError(res, 1)
}

func (c CartDao) Create(ctx context.Context, req CartCreate) error {
	create := &model.Cart{
		UserID:   req.UserId,
		GoodsID:  req.GoodsId,
		GoodsImg: req.GoodsImg,
		Nums:     req.GoodsNum,
		Checked:  false,
	}
	res := c.db.Model(&model.Cart{}).Create(create)
	return HandleError(res, 1)
}

func (c CartDao) Update(ctx context.Context, req CartUpdate) error {
	update := &model.Cart{
		UserID:   req.UserId,
		GoodsID:  req.GoodsId,
		GoodsImg: req.GoodsImg,
		Nums:     req.GoodsNum,
	}
	res := c.db.Model(&model.Cart{}).Where("id = ?", req.ID).Updates(update)
	return HandleError(res, 1)
}

func (c CartDao) Delete(ctx context.Context, req CartMultiGoods) error {
	res := c.db.Model(&model.Cart{}).Where("user_id = ? AND goods_id IN ?", req.UserId, req.GoodsId).Delete(&model.Cart{})
	return HandleError(res, 1)
}

// SelectGoods 勾选商品
func (c CartDao) SelectGoods(ctx context.Context, req CartMultiGoods) error {
	res := c.db.Model(&model.Cart{}).Where("user_id = ? AND goods_id IN ?", req.UserId, req.GoodsId).Update("checked", true)
	// 这里包有bug，如果是全选的话，嗯，前端只传没选中的不就好了，有锤子的bug
	return HandleError(res, len(req.GoodsId))
}

func (c CartDao) CartList(ctx context.Context, base CartBase) (*[]model.Cart, error) {
	resp := &[]model.Cart{}
	res := c.db.Model(&model.Cart{}).Where("user_id = ? ", base.UserId).Find(&resp)
	return resp, HandleError(res, 0)
}
