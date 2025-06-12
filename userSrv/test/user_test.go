package test

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"srv/userSrv/proto"
	"testing"
	"time"
)

var userClient proto.UserClient
var c *grpc.ClientConn

func conn() {
	var err error
	c, err = grpc.NewClient("127.0.0.1:50001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	userClient = proto.NewUserClient(c)
}

func connClose() {
	c.Close()
}

func TestUserList(t *testing.T) {
	conn()
	defer connClose()
	// 测试查询用户列表。并且测试密码验证
	list, err := userClient.GetUserList(context.Background(), &proto.PaginateInfo{
		Page: 1,
		Size: 2,
	})
	if err != nil {
		log.Println(err)
	}

	log.Println("user count:", list.Total)

	for _, user := range list.Data {
		log.Println(user.NickName, user.Password, user.GetPassword())
		ok, err := userClient.CheckPassword(context.Background(), &proto.CheckPasswordRequest{
			Password:        "123456",
			EncryptPassword: user.GetPassword(),
		})
		if err != nil {
			return
		}
		if !ok.Success {
			panic("check password fail")
		}
	}
}

func TestGetUserById(t *testing.T) {
	conn()
	defer connClose()
	user, err := userClient.GetUserById(context.Background(), &proto.IdRequest{Id: 1})
	if err != nil {
		log.Panic(err)
	}

	log.Println(user.NickName, user.Password)
}

func TestGetUserByMobile(t *testing.T) {
	conn()
	defer connClose()
	user, err := userClient.GetUserByMobile(context.Background(), &proto.MobileRequest{Mobile: "17623240006"})
	if err != nil {
		log.Panic(err)
	}

	log.Println(user.NickName, user.Password)
}

//func TestCreateUser(t *testing.T) {
//	conn()
//	defer connClose()
//	i := 10
//	user, err := userClient.CreateUser(context.Background(), &proto.CreateUserRequest{
//		Mobile:   "176232400" + strconv.Itoa(i),
//		Password: "123456",
//		NickName: "batch_add_" + strconv.Itoa(i),
//	})
//	if err != nil {
//		panic(err)
//	}
//	log.Println(user.NickName, user.Password)
//}

func TestUpdateUser(t *testing.T) {
	conn()
	defer connClose()
	_, err := userClient.UpdateUser(context.Background(), &proto.UpdateUserRequest{
		Id:       1,
		NickName: "update_batch_add1xx",
		Birthday: uint64(time.Date(1999, time.July, 1, 0, 0, 0, 0, time.UTC).Unix()),
		Gender:   "male",
	})
	if err != nil {
		panic(err)
	}
}
