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

type CategoryServer struct {
	proto.UnimplementedCategoryServer
}

func (c CategoryServer) CreateCategory(ctx context.Context, info *proto.CreateCategoryInfo) (*proto.CategoryInfoResponse, error) {
	// 分类不能重名是吧
	res := global.DB.Model(&model.Category{}).Where("name = ?", info.Name).First(&model.Category{})
	if res.RowsAffected > 0 {
		return nil, status.Error(codes.InvalidArgument, "分类已经存在")
	}

	// 否则不检测
	if info.ParentCategoryID > 0 {
		res = global.DB.Model(&model.Category{}).Where("parent_category_id = ?", info.ParentCategoryID).First(&model.Category{})
		if res.RowsAffected == 0 {
			return nil, status.Error(codes.InvalidArgument, "待创建的分类的父级不存在")
		}
	}

	category := &model.Category{
		Name:             info.Name,
		ParentCategoryID: info.ParentCategoryID,
		Level:            info.Level,
		IsTab:            false,
	}
	res = global.DB.Model(&model.Category{}).Create(category)
	if res.RowsAffected == 0 {
		return nil, status.Error(codes.InvalidArgument, "创建分类失败")
	}

	return &proto.CategoryInfoResponse{
		ID:               int32(category.ID),
		Name:             info.Name,
		ParentCategoryID: info.ParentCategoryID,
		Level:            info.Level,
		IsTab:            info.IsTab,
	}, nil
}

func (c CategoryServer) DeleteCategory(ctx context.Context, info *proto.DeleteCategoryInfo) (*emptypb.Empty, error) {
	res := global.DB.Model(&model.Category{}).Where("id = ?", info.Id).Delete(&model.Category{})
	if res.RowsAffected == 0 {
		return nil, status.Error(codes.InvalidArgument, "删除分类失败")
	}
	return &emptypb.Empty{}, nil
}

func (c CategoryServer) UpdateCategory(ctx context.Context, info *proto.UpdateCategoryInfo) (*emptypb.Empty, error) {
	res := global.DB.Model(&model.Category{}).Where("id = ?", info.ID).Updates(&model.Category{
		Name:             info.Name,
		ParentCategoryID: info.ParentCategoryID,
		Level:            info.Level,
		IsTab:            info.IsTab,
	})

	if res.RowsAffected == 0 {
		return nil, status.Error(codes.InvalidArgument, "参数错误")
	}
	return &emptypb.Empty{}, nil
}

func (c CategoryServer) GetSubCategory(ctx context.Context, request *proto.CategoryInfoRequest) (*proto.SubCategoryListResponse, error) {
	var resp = &proto.SubCategoryListResponse{}
	var categories []*model.Category
	// 查询当前level和他下面的所有
	res := global.DB.Model(&model.Category{}).Where("level >= ?", request.Level).Find(&categories)
	if res.Error != nil || res.RowsAffected == 0 {
		return resp, status.Error(codes.NotFound, "Category Not Found")
	}
	data := make([]*proto.CategoryInfoResponse, 0, len(categories))
	for _, category := range categories {
		data = append(data, &proto.CategoryInfoResponse{
			ID:               int32(category.ID),
			Name:             category.Name,
			ParentCategoryID: category.ParentCategoryID,
			Level:            category.Level,
			IsTab:            category.IsTab,
		})
	}
	resp.Data = data
	resp.Total = int32(len(categories))
	return resp, nil
}

func (c CategoryServer) GetAllCategories(ctx context.Context, empty *emptypb.Empty) (*proto.CategoryListResponse, error) {
	var resp = &proto.CategoryListResponse{}
	var categories []model.Category
	res := global.DB.Model(&model.Category{}).Find(&categories)
	if res.Error != nil || res.RowsAffected == 0 {
		return resp, status.Error(codes.NotFound, "Category Not Found")
	}
	data := make([]*proto.CategoryInfoResponse, 0, len(categories))
	for _, category := range categories {
		data = append(data, &proto.CategoryInfoResponse{
			ID:               int32(category.ID),
			Name:             category.Name,
			ParentCategoryID: category.ParentCategoryID,
			Level:            category.Level,
			IsTab:            category.IsTab,
		})
	}

	resp.Data = data
	resp.Total = int32(len(categories))
	return resp, nil
}

func (c CategoryServer) mustEmbedUnimplementedCategoryServer() {
	//TODO implement me
	panic("implement me")
}
