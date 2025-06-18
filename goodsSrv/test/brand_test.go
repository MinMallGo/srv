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
	c, err = grpc.NewClient("127.0.0.1:54056", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	brand = proto.NewBrandClient(c)
}

func connClose() {
	c.Close()
}

func TestBrandList(t *testing.T) {
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

func TestBrandCreate(t *testing.T) {
	conn()
	defer connClose()
	res, err := brand.CreateBrand(context.Background(), &proto.CreateBrandInfo{
		Name: "创建测试咯",
		Logo: "xxx",
	})
	if err != nil {
		log.Panicln("CreateBrand : ", err)
	}
	log.Println(res)

	// 这里直接update和delete
	_, err = brand.UpdateBrand(context.Background(), &proto.UpdateBrandInfo{
		ID:   res.ID,
		Name: "改",
	})
	if err != nil {
		log.Panicln("UpdateBrand : ", err)
	}

	_, err = brand.DeleteBrand(context.Background(), &proto.DeleteBrandInfo{
		Id: res.ID,
	})
	if err != nil {
		log.Panicln("DeleteBrand : ", err)
	}
}
