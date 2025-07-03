package model

type Inventory struct {
	BaseID
	GoodsID int32 `json:"goods_id" gorm:"column:goods_id;index:idx_goods"`
	Stocks  int32 `json:"stocks" gorm:"column:stocks"`
	Version int32 `json:"version" gorm:"column:version;default:0"`
	BaseModel
}

type OrderHistory struct {
	BaseID
	OrderSN string       `json:"order_sn" gorm:"column:order_sn;type:varchar(128);index:idx_order,unique"`
	Status  int32        `json:"status" gorm:"column:status;default:1;comment:'1 表示已扣减 2 表示已归还'"`
	Details GoodsDetails `json:"details" gorm:"column:details;type:varchar(255);comment:'记录商品id和数量'"`
	BaseModel
}
