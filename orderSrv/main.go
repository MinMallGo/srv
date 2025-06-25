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
	"srv/orderSrv/global"
	"srv/orderSrv/handler"
	"srv/orderSrv/initialize"
	"srv/orderSrv/proto/gen"
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
	initialize.InitRedLock() // 初始化分布式锁
	flag.Parse()
	if *port == 0 {
		*port = global.GetPort()
	}
	zap.S().Info("ip:%s port:%d\n", *ip, *port)

	// 优雅地注册到注册中心
	//rc := register.NewConsulRegistry(global.SrvConfig.Consul.Host, global.SrvConfig.Consul.Port)
	//id := uuid.New().String()
	//if err := rc.Register(&register.SrvRegisterArgs{
	//	Name: global.SrvConfig.Name,
	//	ID:   id,
	//	Host: global.SrvConfig.Consul.Host,
	//	Port: global.SrvConfig.Consul.Port,
	//	Tags: global.SrvConfig.Consul.Tags,
	//}); err != nil {
	//	zap.L().Panic("register service failed", zap.Error(err))
	//}
	uid := global.UUID()
	client := initialize.RegisterConsul(&initialize.RegArgs{
		Name:    global.SrvConfig.Name,
		ID:      uid,
		Address: "192.168.3.5",
		Port:    *port,
		Tags:    []string{"xxx"},
	})
	// 注册结束
	go func() {
		server := grpc.NewServer(grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(func(p any) (err error) {
				zap.L().Error("recovered from panic", zap.Any("panic", p), zap.Any("stack", string(debug.Stack())))
				return status.Errorf(codes.Internal, "%s", p)
			})),
		))
		// 注册服务
		proto.RegisterOrderServer(server, new(handler.OrderService))
		proto.RegisterCartServer(server, new(handler.CartService))
		//
		lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *ip, *port))
		if err != nil {
			panic(err)
		}
		// 添加内置的health检查
		healthServer := health.NewServer()
		// 这个是检查grpc框架健康的
		grpc_health_v1.RegisterHealthServer(server, healthServer)
		panic(server.Serve(lis))
	}()

	// 添加一个优雅退出
	//quit := make(chan os.Signal)
	//signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	//<-quit
	//err := rc.Deregister(id)
	//if err != nil {
	//	zap.L().Error("deregister service failed", zap.Error(err))
	//}
	//zap.L().Info("deregister service success", zap.String("id", id))
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	err := client.Agent().ServiceDeregister(uid)
	if err != nil {
		zap.L().Error("deregister service failed", zap.Error(err))
	}
	zap.L().Info("deregister service success", zap.String("uid", uid))
}
