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

type CategoryBrandServer struct {
	proto.UnimplementedCategoryBrandServer
}

type CategoryBrandModel struct {
	model.CategoryBrand
	Category *model.Category `gorm:"foreignkey:category_id;references:ID"`
	Brand    *model.Brand    `gorm:"foreignkey:brand_id;references:ID"`
}

// 检查category_id 是否存在
// 检查 brand_id 是否存在

func (c CategoryBrandServer) CreateCategoryBrand(ctx context.Context, info *proto.CreateCategoryBrandInfo) (*proto.CategoryBrandResponse, error) {
	if !CategoryExists(info.CategoryId) {
		return nil, status.Error(codes.InvalidArgument, "分类不存在")
	}

	if !BrandExists(info.BrandId) {
		return nil, status.Error(codes.InvalidArgument, "品牌不存在")
	}

	data := &model.CategoryBrand{
		CategoryID: info.CategoryId,
		BrandID:    info.BrandId,
	}
	res := global.DB.Model(&model.CategoryBrand{}).Create(data)
	if res.RowsAffected == 0 {
		return nil, status.Error(codes.InvalidArgument, "创建品牌分类失败:"+res.Error.Error())
	}
	return &proto.CategoryBrandResponse{
		Id:         int32(data.ID),
		CategoryId: data.CategoryID,
		BrandId:    data.BrandID,
	}, nil
}

func (c CategoryBrandServer) DeleteCategoryBrand(ctx context.Context, info *proto.DeleteCategoryBrandInfo) (*emptypb.Empty, error) {
	res := global.DB.Model(&model.CategoryBrand{}).Where("id = ?", info.Id).Delete(&model.CategoryBrand{})
	if res.RowsAffected == 0 {
		return nil, status.Error(codes.InvalidArgument, "删除品牌分类失败")
	}
	return &emptypb.Empty{}, nil
}

func (c CategoryBrandServer) UpdateCategoryBrand(ctx context.Context, info *proto.UpdateCategoryBrandInfo) (*emptypb.Empty, error) {
	if info.CategoryId != 0 && !CategoryExists(info.CategoryId) {
		return nil, status.Error(codes.InvalidArgument, "分类不存在")
	}

	if info.BrandId != 0 && !BrandExists(info.BrandId) {
		return nil, status.Error(codes.InvalidArgument, "品牌不存在")
	}
	data := &model.CategoryBrand{
		CategoryID: info.CategoryId,
		BrandID:    info.BrandId,
	}
	res := global.DB.Model(&model.CategoryBrand{}).Where("id = ?", info.Id).Updates(data)
	if res.RowsAffected == 0 {
		return nil, status.Error(codes.InvalidArgument, "更新失败")
	}
	return &emptypb.Empty{}, nil
}

func (c CategoryBrandServer) CategoryBrandList(ctx context.Context, request *proto.CategoryBrandInfoRequest) (*proto.CategoryBrandListResponse, error) {
	var resp = &proto.CategoryBrandListResponse{}
	var datas []CategoryBrandModel
	q := global.DB.Model(&model.CategoryBrand{}).Preload("Category").Preload("Brand")

	if request.Id > 0 {
		q = q.Where("id = ?", request.Id)
	}

	if request.CategoryId > 0 {
		q = q.Where("category_id = ?", request.CategoryId)
	}

	if request.BrandId > 0 {
		q = q.Where("brand_id = ?", request.BrandId)
	}
	res := q.Find(&datas)
	if res.RowsAffected == 0 {
		return resp, status.Error(codes.NotFound, "数据不存在")
	}

	dataList := make([]*proto.CategoryBrandResponse, 0, len(datas))
	for _, data := range datas {
		category := &proto.CategoryInfoResponse{
			ID:               int32(data.Category.ID),
			Name:             data.Category.Name,
			ParentCategoryID: data.Category.ParentCategoryID,
			Level:            data.Category.Level,
			IsTab:            data.Category.IsTab,
		}
		brand := &proto.BrandInfoResponse{
			ID:   int32(data.Brand.ID),
			Name: data.Brand.Name,
			Logo: data.Brand.Logo,
		}
		dataList = append(dataList, &proto.CategoryBrandResponse{
			Id:         int32(data.ID),
			CategoryId: data.CategoryID,
			BrandId:    data.BrandID,
			Brand:      brand,
			Category:   category,
		})
	}
	resp.Total = int32(len(datas))
	resp.Data = dataList

	return resp, nil
}

func (c CategoryBrandServer) GetCategoryBrand(ctx context.Context, request *proto.CategoryBrandRequest) (*proto.CategoryBrandResponse, error) {
	var resp = &proto.CategoryBrandResponse{}
	var data CategoryBrandModel
	q := global.DB.Model(&model.CategoryBrand{}).Preload("Category").Preload("Brand").Where("id = ?", request.Id)

	// 先注释掉，感觉没用，我这里就查询一条数据而已
	//if request.CategoryId > 0 {
	//	q = q.Where("category_id = ?", request.CategoryId)
	//}
	//
	//if request.BrandId > 0 {
	//	q = q.Where("brand_id = ?", request.BrandId)
	//}

	res := q.First(&data)
	if res.RowsAffected == 0 {
		return resp, status.Error(codes.NotFound, "数据不存在")
	}

	resp.Category = &proto.CategoryInfoResponse{
		ID:               int32(data.Category.ID),
		Name:             data.Category.Name,
		ParentCategoryID: data.Category.ParentCategoryID,
		Level:            data.Category.Level,
		IsTab:            data.Category.IsTab,
	}

	resp.Brand = &proto.BrandInfoResponse{
		ID:   int32(data.Brand.ID),
		Name: data.Brand.Name,
		Logo: data.Brand.Logo,
	}
	resp.CategoryId = int32(data.Category.ID)
	resp.BrandId = int32(data.Brand.ID)
	return resp, nil
}

func (c CategoryBrandServer) mustEmbedUnimplementedCategoryBrandServer() {
	//TODO implement me
	panic("implement me")
}
