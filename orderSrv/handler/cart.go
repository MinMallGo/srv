package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	dao2 "srv/orderSrv/dao"
	"srv/orderSrv/global"
	proto "srv/orderSrv/proto/gen"
)

type CartService struct {
	proto.UnimplementedCartServer
}

func (c CartService) AddGoods(ctx context.Context, req *proto.AddGoodsReq) (*emptypb.Empty, error) {
	dao := dao2.NewCartDao(global.DB)
	res, err := dao.Get(ctx, dao2.CartBase{UserId: req.UserId, GoodsId: req.GoodsId})
	if err != nil {
		return nil, err
	}

	if res.ID == 0 {
		// 新增
		if err = dao.Create(ctx, dao2.CartCreate{
			UserId:   req.UserId,
			GoodsId:  req.GoodsId,
			GoodsNum: req.GoodsNum,
			GoodsImg: req.GoodsImg,
		}); err != nil {
			return nil, err
		}
		return &emptypb.Empty{}, nil
	}

	if err = dao.Update(ctx, dao2.CartUpdate{
		ID:       int32(res.ID),
		UserId:   res.UserID,
		GoodsId:  res.GoodsID,
		GoodsNum: res.Nums + req.GoodsNum,
		GoodsImg: req.GoodsImg,
	}); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (c CartService) RemoveGoods(ctx context.Context, req *proto.RemoveGoodsReq) (*emptypb.Empty, error) {
	if err := dao2.NewCartDao(global.DB).Delete(ctx, dao2.CartMultiGoods{UserId: req.UserId, GoodsId: req.GoodsId}); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (c CartService) SelectGoods(ctx context.Context, req *proto.SelectGoodsReq) (*emptypb.Empty, error) {
	dao := dao2.NewCartDao(global.DB)
	if !dao.MultiExists(ctx, dao2.CartMultiGoods{UserId: req.UserId, GoodsId: req.GoodsId}) {
		return nil, status.Errorf(codes.NotFound, "请选择正确的商品")
	}

	if err := dao.SelectGoods(ctx, dao2.CartMultiGoods{UserId: req.UserId, GoodsId: req.GoodsId}); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (c CartService) GetCartList(ctx context.Context, req *proto.GetCartListReq) (*proto.CartListResp, error) {
	list, err := dao2.NewCartDao(global.DB).CartList(ctx, dao2.CartBase{UserId: req.UserId})
	if err != nil {
		return nil, err
	}

	data := make([]*proto.CartDetail, 0, len(*list))
	for _, cart := range *list {
		data = append(data, &proto.CartDetail{
			UserID:   cart.UserID,
			GoodsID:  cart.GoodsID,
			GoodsImg: cart.GoodsImg,
			Nums:     cart.Nums,
			Checked:  cart.Checked,
		})
	}

	return &proto.CartListResp{
		Total: int32(len(*list)),
		Data:  data,
	}, nil
}

func (c CartService) UpdateGoodsNum(ctx context.Context, req *proto.UpdateNumReq) (*emptypb.Empty, error) {
	dao := dao2.NewCartDao(global.DB)
	res, err := dao.Get(ctx, dao2.CartBase{UserId: req.UserId, GoodsId: req.GoodsId})
	if err != nil {
		return nil, err
	}

	if err = dao.Update(ctx, dao2.CartUpdate{
		ID:       int32(res.ID),
		UserId:   res.UserID,
		GoodsId:  res.GoodsID,
		GoodsNum: res.Nums + req.GoodsNum,
	}); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (c CartService) mustEmbedUnimplementedCartServer() {
	//TODO implement me
	panic("implement me")
}
