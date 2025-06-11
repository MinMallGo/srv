package model

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt
	IsDeleted bool
}

type User struct {
	ID       int64      `gorm:"primary_key;AUTO_INCREMENT;type:20;"`
	Mobile   string     `gorm:"index:idx_mobile_mobile;unique;type:varchar(11);not null"`
	Password string     `gorm:"type:varchar(120);not null"`
	NickName string     `gorm:"type:varchar(100);not null"`
	Birthday *time.Time `gorm:"column:birthday;type:datetime"`
	Gender   string     `gorm:"column:gender;default:'male';type:varchar(6)"`
	Role     int        `gorm:"column:role;type:5;default:1;comment:'1表示普通用户'"`
	BaseModel
}
