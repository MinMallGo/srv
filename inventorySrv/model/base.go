package model

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt
	IsDeleted bool
}

type BaseID struct {
	ID int64 `gorm:"primary_key;AUTO_INCREMENT;type:int(11);not null;"`
}

type GormList []string

func (g *GormList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), g)
}

func (g GormList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

type Details struct {
	GoodsID  int32
	GoodsNum int32
}

type GoodsDetails []Details

func (g *GoodsDetails) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), g)
}

func (g GoodsDetails) Value() (driver.Value, error) {
	return json.Marshal(g)
}
