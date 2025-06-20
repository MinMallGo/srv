package handler

import (
	"gorm.io/gorm"
	"srv/goodsSrv/global"
	"srv/goodsSrv/model"
)

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func BrandExists(brandId int32) bool {
	res := global.DB.Model(&model.Brand{}).Where("id = ?", brandId).First(&model.Brand{})
	return res.RowsAffected > 0
}

func CategoryExists(categoryId int32) bool {
	res := global.DB.Model(&model.Category{}).Where("id = ?", categoryId).First(&model.Category{})
	return res.RowsAffected > 0
}
