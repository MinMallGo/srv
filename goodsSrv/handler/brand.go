package handler

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"srv/goodsSrv/global"
	"srv/goodsSrv/model"
	proto "srv/goodsSrv/proto/gen"
)

type BrandServer struct {
	proto.UnimplementedBrandServer
}

func (b BrandServer) CreateBrand(ctx context.Context, info *proto.CreateBrandInfo) (*proto.BrandInfoResponse, error) {
	// 1. 检查是否存在，存在则返回错误信息
	res := global.DB.Model(&model.Brand{}).Where("name = ?", info.Name).First(&model.Brand{})

	if res.RowsAffected != 0 {
		zap.L().Error("[BrandServer].CreateBrand", zap.Error(res.Error), zap.String("name", info.Name), zap.String("logo", info.Logo))
		return nil, status.Error(codes.InvalidArgument, "品牌已经存在")
	}

	brand := &model.Brand{
		Name: info.Name,
		Logo: info.Logo,
	}
	res = global.DB.Model(&model.Brand{}).Create(brand)
	if res.RowsAffected == 0 {
		zap.L().Error("[BrandServer].CreateBrand", zap.Error(res.Error), zap.String("name", info.Name), zap.String("logo", info.Logo))
		return nil, status.Error(codes.InvalidArgument, "品牌创建失败")
	}

	return &proto.BrandInfoResponse{
		ID: int32(brand.ID),
	}, nil
}

func (b BrandServer) DeleteBrand(ctx context.Context, info *proto.DeleteBrandInfo) (*emptypb.Empty, error) {
	res := global.DB.Model(&model.Brand{}).Where("id = ?", info.Id).Delete(&model.Brand{})

	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		zap.L().Error("[BrandServer].DeleteBrand 删除失败", zap.Error(res.Error), zap.Int32("id", info.Id))
		return nil, status.Error(codes.Internal, "内部错误")
	}

	if res.RowsAffected == 0 {
		zap.L().Error("[BrandServer].DeleteBrand 删除失败", zap.Error(res.Error), zap.Int32("id", info.Id))
		return nil, status.Error(codes.InvalidArgument, "删除失败")
	}
	return &emptypb.Empty{}, nil
}

func (b BrandServer) UpdateBrand(ctx context.Context, info *proto.UpdateBrandInfo) (*emptypb.Empty, error) {
	res := global.DB.Model(&model.Brand{}).Where("id = ?", info.ID).Updates(&model.Brand{
		Name: info.Name,
		Logo: info.Logo,
	})

	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		zap.L().Error("[BrandServer].UpdateBrand 更新失败", zap.Error(res.Error), zap.String("name", info.Name), zap.String("logo", info.Logo))
		return nil, status.Error(codes.Internal, "内部错误")
	}

	if res.RowsAffected == 0 {
		zap.L().Error("[BrandServer].UpdateBrand 更新失败", zap.Error(res.Error), zap.String("name", info.Name), zap.String("logo", info.Logo))
		return &emptypb.Empty{}, status.Error(codes.InvalidArgument, "<UNK>")
	}
	return &emptypb.Empty{}, nil
}

func (b BrandServer) BrandList(ctx context.Context, request *proto.BrandInfoRequest) (*proto.BrandListResponse, error) {
	var res *proto.BrandListResponse = &proto.BrandListResponse{}
	var count int64

	result := global.DB.Model(&model.Brand{}).Count(&count)
	if result.Error != nil {
		return res, status.Error(codes.Internal, result.Error.Error())
	}

	var brands []*model.Brand
	result = global.DB.Model(&model.Brand{}).Scopes(Paginate(int(request.Page), int(request.PageSize))).Find(&brands)
	if result.Error != nil {
		return res, status.Error(codes.Internal, result.Error.Error())
	}
	brandsResp := make([]*proto.BrandInfoResponse, 0, len(brands))
	for _, brand := range brands {
		brandsResp = append(brandsResp, &proto.BrandInfoResponse{
			ID:   int32(brand.ID),
			Name: brand.Name,
			Logo: brand.Logo,
		})
	}

	res.Data = brandsResp
	res.Total = int32(count)
	return res, nil
}

func (b BrandServer) mustEmbedUnimplementedBrandServer() {
	//TODO implement me
	panic("implement me")
}
