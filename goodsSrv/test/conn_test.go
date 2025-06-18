package test

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	SrvConn *grpc.ClientConn
)

func SrvInit() *grpc.ClientConn {
	var err error
	SrvConn, err = grpc.NewClient(
		fmt.Sprintf(`consul://%s:%d/%s?wait=14s`, "192.168.3.5", 8500, "goods-service"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.L().Fatal("【InitGoodsSrv】服务获取失败：", zap.Error(err))
	}

	return SrvConn
}

func SrvClose() {
	SrvConn.Close()
}
