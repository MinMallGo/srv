package handler

import (
	"context"
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/crypto/pbkdf2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"hash"
	"srv/userSrv/global"
	"srv/userSrv/model"
	"srv/userSrv/proto"
	"strings"
	"time"
)

var (
	DefaultSaltLen   = 10
	DefaultIter      = 100
	DefaultKeyLen    = 32
	DefaultCryptFunc = sha512.New
)

type UserServer struct {
	proto.UnimplementedUserServer
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func User2Response(user model.User) *proto.UserInfoResponse {
	resp := &proto.UserInfoResponse{
		Id:       uint64(user.ID),
		Mobile:   user.Mobile,
		Password: user.Password,
		NickName: user.NickName,
		Gender:   user.Gender,
		Role:     uint32(user.Role),
	}
	if user.Birthday != nil {
		resp.Birthday = uint64(user.Birthday.Unix())
	}
	return resp
}

type passwordEncrypt struct {
	saltLen int
	iter    int
	keyLen  int
	h       func() hash.Hash
}

// pbkdf2
func pwdEncrypt(pwd string) string {
	opts := &passwordEncrypt{
		saltLen: DefaultSaltLen,
		iter:    DefaultIter,
		keyLen:  DefaultKeyLen,
		h:       DefaultCryptFunc,
	}

	// $algorithm$salt$password
	salt, encodePsw := gen(pwd, opts)
	encodePsw = fmt.Sprintf("$pbkdf2-sha256$%s$%s", salt, encodePsw)
	return encodePsw
}

func gen(str string, opts *passwordEncrypt) (string, string) {
	salt := genSalt(opts.saltLen)
	encodePsw := pbkdf2.Key([]byte(str), []byte(salt), opts.iter, opts.keyLen, opts.h)
	return salt, hex.EncodeToString(encodePsw)
}

func verify(password string, salt, encodePwd string, opts *passwordEncrypt) bool {
	tmp := pbkdf2.Key([]byte(password), []byte(salt), opts.iter, opts.keyLen, opts.h)
	return hex.EncodeToString(tmp) == encodePwd
}

func genSalt(strlen int) string {
	b := make([]byte, strlen)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func (u *UserServer) GetUserList(ctx context.Context, info *proto.PaginateInfo) (*proto.UserListResponse, error) {
	resp := &proto.UserListResponse{}
	var userList []model.User
	var count int64
	result := global.DB.Model(&model.User{}).Count(&count)
	if result.Error != nil {
		return resp, status.Errorf(codes.Internal, result.Error.Error())
	}
	resp.Total = uint64(count)

	result = global.DB.Model(&model.User{}).Scopes(Paginate(int(info.Page), int(info.Size))).Find(&userList)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	for _, user := range userList {
		tmp := User2Response(user)
		resp.Data = append(resp.Data, tmp)
	}
	zap.L().Info("获取用户列表")
	return resp, nil
}

func (u *UserServer) GetUserByMobile(ctx context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Model(&model.User{}).Where(model.User{Mobile: req.GetMobile()}).First(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	return User2Response(user), nil
}

func (u *UserServer) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Model(&model.User{}).Where(model.User{ID: int64(req.GetId())}).First(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	return User2Response(user), nil
}

func (u *UserServer) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.UserInfoResponse, error) {
	// 先检查手机号是否被注册了
	var user model.User
	result := global.DB.Model(&model.User{}).Where(model.User{Mobile: req.Mobile}).First(&user)
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已注册")
	}

	user.Mobile = req.Mobile
	user.NickName = req.NickName
	user.Password = pwdEncrypt(req.Password)
	result = global.DB.Model(&model.User{}).Create(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return User2Response(user), nil
}

func (u *UserServer) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*emptypb.Empty, error) {
	// 检查用户是否存在
	var user model.User
	result := global.DB.Model(&model.User{}).Where(model.User{ID: int64(req.Id)}).First(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	birthday := time.Unix(int64(req.GetBirthday()), 0)
	user.NickName = req.GetNickName()
	user.Gender = req.GetGender()
	user.Birthday = &birthday
	result = global.DB.Model(&model.User{}).Where("id = ?", req.GetId()).Updates(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	return &emptypb.Empty{}, nil
}

func (u *UserServer) CheckPassword(ctx context.Context, req *proto.CheckPasswordRequest) (*proto.CheckPasswordResponse, error) {
	// 通过传递过来的密码，验证密码是否正确
	opts := &passwordEncrypt{
		saltLen: DefaultSaltLen,
		iter:    DefaultIter,
		keyLen:  DefaultKeyLen,
		h:       DefaultCryptFunc,
	}
	encrypt := strings.Split(req.GetEncryptPassword(), "$")
	if len(encrypt) < 3 {
		return nil, status.Errorf(codes.InvalidArgument, "参数异常")
	}

	ok := verify(req.GetPassword(), encrypt[2], encrypt[3], opts)
	return &proto.CheckPasswordResponse{Success: ok}, nil
}
