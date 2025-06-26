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

type OrderListReq struct {
	UserID int32
	Page   int32
	Size   int32
}

type OrderDetailReq struct {
	UserID  int32
	OrderId int32
	OrderSN string
}

type Order struct {
	model.Order
	Goods []*model.OrderGoods `gorm:"foreignkey:OrderId;references:ID"`
}

func (r *OrderDao) Create(param *model.Order) (int, error) {
	res := r.db.Model(&model.Order{}).Create(&param)
	return int(param.ID), HandleError(res, 1)
}

func (r *OrderDao) GetList(ctx context.Context, req OrderListReq) (*proto.OrderListResp, error) {
	list := []Order{}
	var count int64
	res := r.db.Model(&model.Order{})
	if req.UserID != 0 {
		res = res.Where("user_id = ?", req.UserID)
	}
	res = res.Count(&count)
	err := HandleError(res, 0)
	if err != nil {
		return nil, err
	}

	x := r.db.Model(&model.Order{}).Preload("Goods").Scopes(Paginate(int(req.Page), int(req.Size)))
	if req.UserID != 0 {
		x = x.Where("user_id = ?", req.UserID)
	}
	x = x.Find(&list)
	err = HandleError(x, 0)
	if err != nil {
		return nil, err
	}

	data := make([]*proto.OrderDetailResp, 0, len(list))
	for _, order := range list {
		goods := make([]*proto.GoodsInfo, 0, len(order.Goods))
		for _, good := range order.Goods {
			goods = append(goods, &proto.GoodsInfo{
				OrderID:    good.OrderId,
				OrderSN:    good.OrderSN,
				GoodsID:    good.GoodsId,
				GoodsPrice: good.GoodsPrice,
				PayPrice:   good.PayPrice,
				GoodsName:  good.GoodsName,
				Num:        good.Nums,
			})
		}
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
			CreateAt:        order.CreatedAt.Format("2006-01-02 15:04:05"),
			Goods:           goods,
		})
	}

	return &proto.OrderListResp{
		Total: int32(count),
		Data:  data,
	}, nil
}

func (r *OrderDao) GetDetail(ctx context.Context, req OrderDetailReq) (*Order, error) {
	order := &Order{}
	query := r.db.Model(&model.Order{}).Preload("Goods")
	if req.OrderId > 0 {
		query = query.Where("id = ?", req.OrderId)
	}

	if len(req.OrderSN) > 0 {
		query = query.Where("order_sn = ?", req.OrderSN)
	}

	if req.UserID > 0 {
		query = query.Where("user_id = ?", req.UserID)
	}

	res := query.Find(order)
	err := HandleError(res, 1)
	return order, err
}
