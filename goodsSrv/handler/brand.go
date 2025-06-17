package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"srv/goodsSrv/global"
	"srv/goodsSrv/model"
	proto "srv/goodsSrv/proto/gen"
)

type BrandServer struct {
	proto.UnimplementedBrandServer
}

func (b BrandServer) CreateBrand(ctx context.Context, info *proto.CreateBrandInfo) (*proto.BrandInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (b BrandServer) DeleteBrand(ctx context.Context, info *proto.DeleteBrandInfo) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (b BrandServer) UpdateBrand(ctx context.Context, info *proto.UpdateBrandInfo) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
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
