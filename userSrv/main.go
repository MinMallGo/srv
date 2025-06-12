package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"srv/userSrv/handler"
	"srv/userSrv/proto"
)

func main() {
	// 先通过flag包获取用户的环境量表输入
	ip := flag.String("ip", "0.0.0.0", "ip address")
	port := flag.Int("port", 50001, "port number")
	flag.Parse()
	fmt.Printf("ip:%s port:%d\n", *ip, *port)

	server := grpc.NewServer()
	proto.RegisterUserServer(server, new(handler.UserServer))
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *ip, *port))
	if err != nil {
		panic(err)
	}
	panic(server.Serve(lis))
}
