package global

import (
	"github.com/go-redsync/redsync/v4"
	"gorm.io/gorm"
	"srv/inventorySrv/structs"
)

var SrvConfig = &structs.ServerConfig{}
var DB *gorm.DB
var RedLock *redsync.Redsync = &redsync.Redsync{}
