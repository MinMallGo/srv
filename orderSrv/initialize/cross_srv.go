package initialize

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"srv/orderSrv/global"
	proto "srv/orderSrv/proto/gen"
)

// InitCrossSrv 连接第三方服务
func InitCrossSrv() {
	log.Printf("%#v", global.SrvConfig.CrossSrv)
	ConnGoods()
	ConnStock()
}

// ConnGoods goods-service
func ConnGoods() {
	conn, err := grpc.NewClient(
		fmt.Sprintf(`consul://%s:%d/%s?wait=14s`, "192.168.3.5", 8500, global.SrvConfig.CrossSrv.GoodsSrv),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.L().Fatal(fmt.Sprintf("【%s】服务获取失败：", global.SrvConfig.CrossSrv.GoodsSrv), zap.Error(err))
	}

	global.CrossSrv.Goods = proto.NewGoodsClient(conn)
}

// ConnStock stock-service
func ConnStock() {
	conn, err := grpc.NewClient(
		fmt.Sprintf(`consul://%s:%d/%s?wait=14s`, "192.168.3.5", 8500, global.SrvConfig.CrossSrv.Inventory),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.L().Fatal(fmt.Sprintf("【%s】服务获取失败：", global.SrvConfig.CrossSrv.Inventory), zap.Error(err))
	}

	global.CrossSrv.Inventory = proto.NewInventoryClient(conn)
}
