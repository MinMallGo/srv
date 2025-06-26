package dao

import (
	"context"
	"gorm.io/gorm"
	"srv/orderSrv/model"
	proto "srv/orderSrv/proto/gen"
)

type OrderDao struct {
	db *gorm.DB
}

func NewOrderDao(db *gorm.DB) *OrderDao {
	return &OrderDao{db: db}
}

type OrderListResp struct {
	Page int32
	Size int32
}

type OrderDetailResp struct {
	OrderId int32
	OrderSN string
}

type Order struct {
	model.Order
	Goods []*model.OrderGoods `gorm:"foreignkey:OrderId;references:ID"`
}

func (r *OrderDao) Create(param model.Order) (int, error) {
	res := r.db.Model(&model.Order{}).Create(&param)
	return int(param.ID), HandleError(res, 1)
}

func (r *OrderDao) GetList(ctx context.Context, req OrderListResp) (*proto.OrderListResp, error) {
	list := []model.Order{}
	var count int64

	res := r.db.Model(&model.Order{}).Count(&count)
	err := HandleError(res, 0)
	if err != nil {
		return nil, err
	}

	res = r.db.Model(&model.Order{}).Scopes(Paginate(int(req.Page), int(req.Size))).Find(&list)
	err = HandleError(res, 0)
	if err != nil {
		return nil, err
	}

	data := make([]*proto.OrderDetailResp, 0, len(list))
	for _, order := range list {
		data = append(data, &proto.OrderDetailResp{
			UserID:          order.UserID,
			OrderSN:         order.OrderSN,
			PayType:         order.PayType,
			Status:          order.Status,
			TradeNo:         order.TradeNo,
			SubjectTitle:    order.SubjectTitle,
			OrderPrice:      order.OrderPrice,
			FinalPrice:      order.FinalPrice,
			Address:         order.Address,
			RecipientName:   order.RecipientName,
			RecipientMobile: order.RecipientMobile,
			Message:         order.Message,
			Snapshot:        order.Snapshot,
		})
	}

	return &proto.OrderListResp{
		Total: int32(count),
		Data:  data,
	}, nil
}

func (r *OrderDao) GetDetail(ctx context.Context, req OrderDetailResp) (*Order, error) {
	order := &Order{}
	query := r.db.Model(&model.Order{}).Preload("Goods")
	if req.OrderId > 0 {
		query = query.Where("id = ?", req.OrderId)
	}

	if len(req.OrderSN) > 0 {
		query = query.Where("order_sn = ?", req.OrderSN)
	}

	res := query.Find(order)
	err := HandleError(res, 1)
	return order, err
}
