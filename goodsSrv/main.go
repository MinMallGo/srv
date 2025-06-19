package main

import (
	"flag"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
	"net"
	"os"
	"os/signal"
	"runtime/debug"
	"srv/goodsSrv/global"
	"srv/goodsSrv/handler"
	"srv/goodsSrv/initialize"
	"srv/goodsSrv/proto/gen"
	"syscall"
)

func main() {
	// 先通过flag包获取用户的环境量表输入
	ip := flag.String("ip", "0.0.0.0", "ip address")
	port := flag.Int("port", 0, "port number")
	/*
		1. 初始化日志
		2. 初始化配置文件
		3. 初始化db
	*/
	initialize.InitZap()
	initialize.InitConfig()
	initialize.InitDB()
	flag.Parse()
	if *port == 0 {
		*port = global.GetPort()
	}
	zap.S().Info("ip:%s port:%d\n", *ip, *port)
	uid := global.UUID()
	client := initialize.RegisterConsul(&initialize.RegArgs{
		Name:    global.SrvConfig.Name,
		ID:      uid,
		Address: "192.168.3.5",
		Port:    *port,
		Tags:    []string{"xxx"},
	})

	go func() {
		server := grpc.NewServer(grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(func(p any) (err error) {
				zap.L().Error("recovered from panic", zap.Any("panic", p), zap.Any("stack", string(debug.Stack())))
				return status.Errorf(codes.Internal, "%s", p)
			})),
		))

		proto.RegisterBrandServer(server, new(handler.BrandServer))
		proto.RegisterBannerServer(server, new(handler.BannerServer))
		proto.RegisterCategoryServer(server, new(handler.CategoryServer))
		proto.RegisterCategoryBrandServer(server, new(handler.CategoryBrandServer))
		proto.RegisterGoodsServer(server, new(handler.GoodsServer))

		//
		lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *ip, *port))
		if err != nil {
			panic(err)
		}
		// 添加内置的health检查
		healthServer := health.NewServer()
		grpc_health_v1.RegisterHealthServer(server, healthServer)
		panic(server.Serve(lis))
	}()

	// 添加一个优雅退出
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	err := client.Agent().ServiceDeregister(uid)
	if err != nil {
		zap.L().Error("deregister service failed", zap.Error(err))
	}
	zap.L().Info("deregister service success", zap.String("uid", uid))
}
