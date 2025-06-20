package model

type Category struct {
	BaseID
	Name             string `gorm:"type:varchar(100);not null"`
	ParentCategoryID int32
	//ParentCategory   *Category
	Level int32 `gorm:"type:int;default:1;not null"`
	IsTab bool  `gorm:"default:false;not null"`
	BaseModel
}

type Brand struct {
	BaseID
	Name string `gorm:"type:varchar(100);not null"`
	Logo string `gorm:"type:varchar(100);not null"`
	BaseModel
}

type CategoryBrand struct {
	BaseID
	CategoryID int32 `gorm:"type:int(11);index:idx_category_brands,unique"`
	//Category   *Category
	BrandID int32 `gorm:"type:int(11);index:idx_category_brands,unique"`
	//Brand      *Brand
	BaseModel
}

type Banner struct {
	BaseID
	Image string `gorm:"type:varchar(100);not null"`
	Url   string `gorm:"type:varchar(100);not null"`
	Index int32  `gorm:"type:int(11);default:1;not null"`
	BaseModel
}

type Goods struct {
	BaseID
	CategoryID int32 `gorm:"type:int(11);not null"`
	//Category   *Category
	BrandID int32 `gorm:"type:int(11);not null"`
	//Brand      *Brand

	OnSale   *bool `gorm:"default:false;not null"`
	IsNew    *bool `gorm:"default:false;not null"`
	IsHot    *bool `gorm:"default:false;not null"`
	ShipFree bool  `gorm:"default:false;not null"`

	Name            string   `gorm:"type:varchar(100);not null"`
	GoodsSn         string   `gorm:"type:varchar(100);not null"`
	ClickNum        int32    `gorm:"type:int(11);default:0;not null"`
	SoldNum         int32    `gorm:"type:int(11);default:0;not null"`
	FavNum          int32    `gorm:"type:int(11);default:0;not null"`
	MarketPrice     float32  `gorm:"default:0;not null"`
	ShopPrice       float32  `gorm:"default:0;not null"`
	GoodsBrief      string   `gorm:"type:varchar(100);not null"`
	ImageUrl        GormList `gorm:"type:varchar(1000);not null"`
	Description     GormList `gorm:"type:varchar(1000);not null"`
	GoodsFrontImage string   `gorm:"type:varchar(200);not null"`
	BaseModel
}
