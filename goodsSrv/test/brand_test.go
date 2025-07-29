package test

import (
	"context"
	"goodsSrv/proto/gen"
	"google.golang.org/grpc"
	"log"
	"testing"
)

var brand proto.BrandClient
var c *grpc.ClientConn

func TestBrandList(t *testing.T) {
	c = SrvInit()
	defer SrvClose()
	brand = proto.NewBrandClient(c)
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
	c = SrvInit()
	defer SrvClose()
	brand = proto.NewBrandClient(c)
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
