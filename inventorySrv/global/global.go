package global

import (
	"gorm.io/gorm"
	"srv/inventorySrv/structs"
)

var SrvConfig = &structs.ServerConfig{}
var DB *gorm.DB
