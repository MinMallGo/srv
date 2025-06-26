package global

import (
	"github.com/go-redsync/redsync/v4"
	"gorm.io/gorm"
	proto "srv/orderSrv/proto/gen"
	"srv/orderSrv/structs"
)

type CrossServer struct {
	Goods     proto.GoodsClient
	Inventory proto.InventoryClient // 库存服务，指的是 inventory
}

var SrvConfig = &structs.ServerConfig{}
var DB *gorm.DB
var RedLock *redsync.Redsync = &redsync.Redsync{}
var CrossSrv = &CrossServer{}
