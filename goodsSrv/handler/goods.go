package handler

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"srv/goodsSrv/global"
	"srv/goodsSrv/model"
	proto "srv/goodsSrv/proto/gen"
)

type GoodsServer struct {
	proto.UnimplementedGoodsServer
}

// GoodsModel2 用于关联数据查询等
type GoodsModel2 struct {
	model.Goods
	Category *model.Category `gorm:"foreignkey:category_id;references:ID"`
	Brand    *model.Brand    `gorm:"foreignkey:brand_id;references:ID"`
}

func (g GoodsServer) GoodsList(ctx context.Context, request *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	var resp = &proto.GoodsListResponse{}
	var goods []*GoodsModel2
	x := global.DB.Model(&model.Goods{}).Preload("Category").Preload("Brand")
	if request.IsTab {
		x = x.Where("is_tab = ?", request.IsTab)
	}
	if request.IsNew {
		x = x.Where("is_new = ?", request.IsNew)
	}
	if request.IsHot {
		x = x.Where("is_hot = ?", request.IsHot)
	}
	if request.KeyWord != "" {
		x = x.Where("name LIKE ?", fmt.Sprintf(`%%%s%%`, request.KeyWord))
	}
	if request.PriceMin > 0 {
		x = x.Where("shop_price >= ?", request.PriceMin)
	}
	if request.PriceMax > 0 {
		x = x.Where("shop_price <= ?", request.PriceMax)
	}

	if request.TopCategory > 0 {
		x = x.Joins("left join category on category.id = goods.category_id").Where("category.level >= ?", request.TopCategory)
	}

	var total int64
	y := x
	res := y.Count(&total)
	if res.RowsAffected == 0 {
		return nil, status.Error(codes.InvalidArgument, "未查询到商品")
	}

	res = x.Scopes(Paginate(int(request.Pages), int(request.PageSize))).Find(&goods)
	if res.RowsAffected == 0 {
		return nil, status.Error(codes.InvalidArgument, "未查询到商品")
	}

	data := g.TidyGoods(goods)

	resp.Total = int32(total)
	resp.Data = data
	return resp, nil
}

// TidyGoods 整理商品
func (g GoodsServer) TidyGoods(goods []*GoodsModel2) []*proto.GoodsInfoResponse {
	data := make([]*proto.GoodsInfoResponse, 0, len(goods))
	for _, good := range goods {
		data = append(data, &proto.GoodsInfoResponse{
			Id:          int32(good.ID),
			CategoryId:  good.CategoryID,
			BrandId:     good.BrandID,
			OnSale:      good.OnSale,
			ShipFree:    good.ShipFree,
			IsNew:       good.IsNew,
			IsHot:       good.IsHot,
			Name:        good.Name,
			GoodsSn:     good.GoodsSn,
			ClickNum:    good.ClickNum,
			SoldNum:     good.SoldNum,
			FavNum:      good.FavNum,
			MarketPrice: good.MarketPrice,
			ShopPrice:   good.ShopPrice,
			GoodsBrief:  good.GoodsBrief,
			//ImageUrl:        good.ImageUrl,
			//Description:     good.Description,
			GoodsFrontImage: good.GoodsFrontImage,
			CreatedAt:       uint64(good.CreatedAt.Unix()),
			IsDeleted:       good.IsDeleted,
			Category: &proto.CategoryBriefInfoResponse{
				Id:   int32(good.Category.ID),
				Name: good.Category.Name,
			},
			Brand: &proto.BrandInfoResponse{
				ID:   int32(good.Brand.ID),
				Name: good.Brand.Name,
				Logo: good.Brand.Logo,
			},
		})
	}

	return data
}

func (g GoodsServer) BatchGetGoods(ctx context.Context, info *proto.BatchGoodsInfo) (*proto.GoodsListResponse, error) {
	var resp = &proto.GoodsListResponse{}
	var goods []*GoodsModel2
	x := global.DB.Model(&model.Goods{}).Preload("Category").Preload("Brand").Where("id in ?", info.Id).Find(&goods)
	if x.RowsAffected == 0 {
		return nil, status.Error(codes.InvalidArgument, "未查询到商品")
	}
	resp.Total = int32(int64(len(goods)))
	resp.Data = g.TidyGoods(goods)
	return resp, nil
}

func (g GoodsServer) BrandExists(brandId int32) bool {
	res := global.DB.Model(&model.Brand{}).Where("id = ?", brandId).First(&model.Brand{})
	return res.RowsAffected > 0
}

func (g GoodsServer) CategoryExists(categoryId int32) bool {
	res := global.DB.Model(&model.Category{}).Where("id = ?", categoryId).First(&model.Category{})
	return res.RowsAffected > 0
}

func (g GoodsServer) CreateGods(ctx context.Context, info *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	if !g.BrandExists(info.BrandId) {
		return nil, status.Error(codes.InvalidArgument, "品牌不存在")
	}
	if !g.CategoryExists(info.CategoryId) {
		return nil, status.Error(codes.InvalidArgument, "分类不存在")
	}

	goods := &model.Goods{
		CategoryID:      info.CategoryId,
		BrandID:         info.BrandId,
		OnSale:          info.OnSale,
		ShipFree:        info.ShipFree,
		IsNew:           info.IsNew,
		Name:            info.Name,
		GoodsSn:         info.GoodsSn,
		ClickNum:        info.ClickNum,
		SoldNum:         info.SoldNum,
		FavNum:          info.FavNum,
		MarketPrice:     info.MarketPrice,
		ShopPrice:       info.ShopPrice,
		GoodsBrief:      info.GoodsBrief,
		ImageUrl:        info.ImageUrl,
		Description:     info.Description,
		GoodsFrontImage: info.GoodsFrontImage,
	}

	res := global.DB.Model(&model.Goods{}).Create(goods)
	if res.RowsAffected == 0 {
		return nil, status.Error(codes.InvalidArgument, "创建商品失败")
	}
	return &proto.GoodsInfoResponse{Id: int32(goods.ID)}, nil
}

func (g GoodsServer) UpdateGoods(ctx context.Context, info *proto.UpdateGoodsInfo) (*emptypb.Empty, error) {
	if !g.BrandExists(info.BrandId) {
		return nil, status.Error(codes.InvalidArgument, "品牌不存在")
	}
	if !g.CategoryExists(info.CategoryId) {
		return nil, status.Error(codes.InvalidArgument, "分类不存在")
	}

	goods := &model.Goods{
		CategoryID:      info.CategoryId,
		BrandID:         info.BrandId,
		OnSale:          info.OnSale,
		ShipFree:        info.ShipFree,
		IsNew:           info.IsNew,
		Name:            info.Name,
		GoodsSn:         info.GoodsSn,
		ClickNum:        info.ClickNum,
		SoldNum:         info.SoldNum,
		FavNum:          info.FavNum,
		MarketPrice:     info.MarketPrice,
		ShopPrice:       info.ShopPrice,
		GoodsBrief:      info.GoodsBrief,
		ImageUrl:        info.ImageUrl,
		Description:     info.Description,
		GoodsFrontImage: info.GoodsFrontImage,
	}
	res := global.DB.Model(&model.Goods{}).Where("id = ?", info.Id).Updates(goods)
	if res.RowsAffected == 0 {
		return nil, status.Error(codes.InvalidArgument, "更新商品失败")
	}
	return &emptypb.Empty{}, nil
}

func (g GoodsServer) DeleteGoods(ctx context.Context, info *proto.DeleteGoodsInfo) (*emptypb.Empty, error) {
	res := global.DB.Model(&model.Goods{}).Where("id = ?", info.Id).Delete(&model.Goods{})
	if res.RowsAffected == 0 {
		return nil, status.Error(codes.InvalidArgument, "删除商品失败")
	}
	return &emptypb.Empty{}, nil
}

func (g GoodsServer) GetGoodsDetail(ctx context.Context, request *proto.GoodsInfoRequest) (*proto.GoodsInfoResponse, error) {
	var goods []*GoodsModel2
	res := global.DB.Model(&model.Goods{}).Preload("Category").Preload("Brand").Where("id = ?", request.Id).Find(&goods)
	if res.RowsAffected == 0 {
		return nil, status.Error(codes.InvalidArgument, "未查询到商品")
	}
	data := g.TidyGoods(goods)
	return data[0], nil
}

func (g GoodsServer) mustEmbedUnimplementedGoodsServer() {
	//TODO implement me
	panic("implement me")
}
