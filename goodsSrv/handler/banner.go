package handler

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"srv/goodsSrv/global"
	"srv/goodsSrv/model"
	proto "srv/goodsSrv/proto/gen"
)

type BannerServer struct {
	proto.UnimplementedBannerServer
}

func (b *BannerServer) CreateBanner(ctx context.Context, info *proto.CreateBannerInfo) (*proto.BannerInfoResponse, error) {
	banner := &model.Banner{
		Image: info.Image,
		Url:   info.Url,
		Index: info.Index,
	}
	res := global.DB.Model(&model.Banner{}).Create(banner)

	if res.RowsAffected == 0 {
		zap.L().Error("create Banner failed", zap.Any("info", info))
		return nil, status.Error(codes.InvalidArgument, "创建轮播图失败")
	}

	return &proto.BannerInfoResponse{
		ID:    int32(banner.ID),
		Image: banner.Image,
		Url:   banner.Url,
		Index: banner.Index,
	}, nil

}

func (b *BannerServer) DeleteBanner(ctx context.Context, info *proto.DeleteBannerInfo) (*emptypb.Empty, error) {
	res := global.DB.Model(&model.Banner{}).Where("id = ?", info.Id).Delete(&model.Banner{})
	if res.RowsAffected == 0 {
		zap.L().Error("delete Banner failed", zap.Any("info", info))
		return nil, status.Error(codes.InvalidArgument, "删除轮播图失败")
	}
	return &emptypb.Empty{}, nil
}

func (b *BannerServer) UpdateBanner(ctx context.Context, info *proto.UpdateBannerInfo) (*emptypb.Empty, error) {
	res := global.DB.Model(&model.Banner{}).Where("id = ?", info.ID).Updates(&model.Banner{
		Image: info.Image,
		Url:   info.Url,
		Index: 0,
	})
	if res.RowsAffected == 0 {
		zap.L().Error("update Banner failed", zap.Any("info", info))
		return nil, status.Error(codes.InvalidArgument, "更新轮播图失败")
	}
	return &emptypb.Empty{}, nil
}

func (b *BannerServer) BannerList(ctx context.Context, request *proto.BannerInfoRequest) (*proto.BannerListResponse, error) {
	var resp *proto.BannerListResponse = &proto.BannerListResponse{}
	banners := make([]model.Banner, 0)

	_ = global.DB.Model(&model.Banner{}).Find(&banners)

	data := make([]*proto.BannerInfoResponse, 0, len(banners))
	for _, banner := range banners {
		data = append(data, &proto.BannerInfoResponse{
			ID:    int32(banner.ID),
			Image: banner.Image,
			Url:   banner.Url,
			Index: banner.Index,
		})
	}

	resp.Total = int32(0)
	resp.Data = data
	return resp, nil
}

func (b *BannerServer) mustEmbedUnimplementedBannerServer() {
	//TODO implement me
	panic("implement me")
}
