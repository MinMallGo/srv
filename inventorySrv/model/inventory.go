package model

type Inventory struct {
	BaseID
	GoodsID int32 `json:"goods_id" gorm:"column:goods_id;index:idx_goods"`
	Stocks  int32 `json:"stocks" gorm:"column:stocks"`
	Version int32 `json:"version" gorm:"column:version;default:0"`
	BaseModel
}
