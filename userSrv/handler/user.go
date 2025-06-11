package handler

import (
	"context"
	"gorm.io/gorm"
	"srv/userSrv/global"
	"srv/userSrv/model"
	"srv/userSrv/proto"
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

func (u UserServer) GetUserList(ctx context.Context, info *proto.PaginateInfo) (*proto.UserListResponse, error) {
	resp := &proto.UserListResponse{}
	var userList []model.User
	var count int64
	result := global.DB.Count(&count)
	if result.Error != nil {
		return resp, result.Error
	}
	resp.Total = uint64(count)

	result = global.DB.Scopes(Paginate(int(info.Page), int(info.Size))).Find(&userList)
	if result.Error != nil {
		return nil, result.Error
	}
	for _, user := range userList {
		tmp := User2Response(user)
		resp.Data = append(resp.Data, tmp)
	}

	return resp, nil
}
