package main

import (
	"crypto/sha512"
	"fmt"
	"goodsSrv/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	DefaultSaltLen   = 10
	DefaultIter      = 100
	DefaultKeyLen    = 32
	DefaultCryptFunc = sha512.New
)

func main() {
	migrate()
	//pwdEncrypt("gen password")
}

func migrate() {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "root:root@tcp(127.0.0.1:3306)/shop_goods_srv?charset=utf8&parseTime=True&loc=Local", // DSN data source name
		DefaultStringSize:         256,                                                                                  // string 类型字段的默认长度
		DisableDatetimePrecision:  true,                                                                                 // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,                                                                                 // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,                                                                                 // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,                                                                                // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(fmt.Sprintf("autoMigrate failed: %v", err))
	}
	err = db.AutoMigrate(&model.Category{}, &model.Banner{}, &model.Goods{}, &model.CategoryBrand{}, &model.Brand{})
	if err != nil {
		panic(fmt.Sprintf("autoMigrate failed: %v", err))
	}
}
