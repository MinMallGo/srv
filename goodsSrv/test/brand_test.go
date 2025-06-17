package test

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"srv/goodsSrv/proto/gen"
	"testing"
)

var brand proto.BrandClient
var c *grpc.ClientConn

func conn() {
	var err error
	c, err = grpc.NewClient("127.0.0.1:50522", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	brand = proto.NewBrandClient(c)
}

func connClose() {
	c.Close()
}

func TestUserList(t *testing.T) {
	conn()
	defer connClose()
	// 测试查询用户列表。并且测试密码验证
	list, err := brand.BrandList(context.Background(), &proto.BrandInfoRequest{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		log.Println(err)
	}

	log.Println("user count:", list.Total)

	for _, b := range list.Data {
		log.Println(b.Name)
	}
}
