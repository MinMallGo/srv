package global

import (
	"gorm.io/gorm"
	"srv/goodsSrv/structs"
)

var SrvConfig = &structs.ServerConfig{}
var DB *gorm.DB
