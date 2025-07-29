package global

import (
	"goodsSrv/structs"
	"gorm.io/gorm"
)

var SrvConfig = &structs.ServerConfig{}
var DB *gorm.DB
