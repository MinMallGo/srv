package main

import (
	"flag"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"srv/userSrv/global"
	"srv/userSrv/handler"
	"srv/userSrv/initialize"
	"srv/userSrv/proto"
)

func main() {
	// 先通过flag包获取用户的环境量表输入
	ip := flag.String("ip", "0.0.0.0", "ip address")
	port := flag.Int("port", 50001, "port number")
	/*
		1. 初始化日志
		2. 初始化配置文件
		3. 初始化db
	*/
	initialize.InitZap()
	initialize.InitConfig()
	initialize.InitDB()
	flag.Parse()
	zap.S().Info("ip:%s port:%d\n", *ip, *port)
	initialize.RegisterConsul(&initialize.RegArgs{
		Name:    global.SrvConfig.Name,
		ID:      global.SrvConfig.Name,
		Address: "192.168.3.5",
		Port:    *port,
		Tags:    []string{"xxx"},
	})

	server := grpc.NewServer()
	proto.RegisterUserServer(server, new(handler.UserServer))
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *ip, *port))
	if err != nil {
		panic(err)
	}
	// 添加内置的health检查
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(server, healthServer)
	panic(server.Serve(lis))
}
