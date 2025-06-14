package global

import (
	"gorm.io/gorm"
	"srv/userSrv/structs"
)

var SrvConfig = &structs.ServerConfig{}
var DB *gorm.DB
