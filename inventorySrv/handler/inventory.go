package handler

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	dao2 "srv/inventorySrv/dao"
	"srv/inventorySrv/global"
	"srv/inventorySrv/model"
	"srv/inventorySrv/proto/gen"
)

type InventoryServer struct {
	proto.UnimplementedInventoryServer
}

func (i InventoryServer) SetStock(ctx context.Context, info *proto.SetInfo) (*emptypb.Empty, error) {
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		// 通过id查询是否存在，不存在则新增，存在则修改
		stock := &model.Inventory{}
		res := tx.Model(&model.Inventory{}).Where("goods_id = ?", info.GoodsId).First(stock)
		if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
			zap.L().Info("<SetStock>.First(stock)", zap.Error(res.Error))
			return status.Error(codes.Internal, res.Error.Error())
		}

		if stock.ID == 0 {
			// 新增
			res = tx.Model(&model.Inventory{}).Create(&model.Inventory{
				GoodsID: info.GoodsId,
				Stocks:  info.Stock,
			})
			if res.Error != nil {
				zap.L().Info("<SetStock>.Create(&model.Inventory", zap.Error(res.Error))
				return status.Error(codes.Internal, res.Error.Error())
			}
			return nil
		}

		// 修改
		res = tx.Model(&model.Inventory{}).Where("goods_id = ?", info.GoodsId).Update("stocks", info.Stock)
		if res.Error != nil {
			zap.L().Info(`<SetStock>.Update("stocks", info.Stock)`, zap.Error(res.Error))
			return status.Error(codes.Internal, res.Error.Error())
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (i InventoryServer) SellStock(ctx context.Context, info *proto.MultipleInfo) (*emptypb.Empty, error) {
	stocks := make([]Stocks, 0, len(info.Sell))
	for _, data := range info.Sell {
		stocks = append(stocks, Stocks{
			GoodsID: data.GoodsId,
			Stock:   data.Stock,
		})
	}

	//// 使用悲观锁来进行更新。test ✅
	//pessimism := func() (*emptypb.Empty, error) {
	//	consistency := GetConsistency(0)
	//	err := consistency.Decr(stocks...)
	//	if err != nil {
	//		return nil, err
	//	}
	//	return &emptypb.Empty{}, nil
	//}
	//// 使用悲观锁来进行更新。test ✅

	//// 使用乐观锁来进行更新
	//// 😒❌ 这里真的要注意，因为数据竞争太大了，如果导致重试的次数过少的话，会大大影响成功率
	//optimism := func() (*emptypb.Empty, error) {
	//	consistency := GetConsistency(1)
	//	err := consistency.Decr(stocks...)
	//	if err != nil {
	//		return nil, err
	//	}
	//	return &emptypb.Empty{}, nil
	//}
	//return optimism()
	// 使用分布式锁进行更新
	redLock := func() (*emptypb.Empty, error) {
		consistency := GetConsistency(99)
		err := consistency.Decr(stocks...)
		if err != nil {
			return nil, err
		}
		return &emptypb.Empty{}, nil
	}
	return redLock()

}

func (i InventoryServer) GetStock(ctx context.Context, info *proto.GetInfo) (*proto.StockResp, error) {
	stock := &model.Inventory{}
	res := global.DB.Model(&model.Inventory{}).Where("goods_id = ?", info.GoodsId).First(stock)
	if res.Error != nil {
		zap.L().Info("<GetStock>.First(stock)", zap.Error(res.Error))
		return nil, status.Error(codes.Internal, res.Error.Error())
	}
	if stock.ID == 0 {
		zap.L().Info(`<GetStock>.StockID == 0`)
		return nil, status.Error(codes.Internal, "商品不存在")
	}

	resp := &proto.StockResp{
		GoodsId: stock.GoodsID,
		Stock:   stock.Stocks,
	}
	return resp, nil

}

func (i InventoryServer) ReturnStock(ctx context.Context, info *proto.MultipleInfo) (*emptypb.Empty, error) {
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		// 用于快速IN查询
		goods := make([]int32, 0, len(info.Sell))
		// 用于调用封装的减库存
		incr := make([]dao2.Stock, 0, len(info.Sell))
		for _, v := range info.Sell {
			goods = append(goods, v.GoodsId)
			incr = append(incr, dao2.Stock{GoodsId: v.GoodsId, Stocks: v.Stock})
		}

		infos := &[]model.Inventory{}
		// 加读锁进行查询
		res := tx.Model(&model.Inventory{}).Clauses(clause.Locking{Strength: "UPDATE"}).Where("goods_id in ?", goods).Find(&infos)
		if res.Error != nil {
			zap.L().Info("<SellStock>.Find(goods)", zap.Error(res.Error))
			return status.Error(codes.Internal, res.Error.Error())
		}

		if res.RowsAffected != int64(len(goods)) {
			zap.L().Info(`<SellStock>.RowsAffected != int64(len(goods))`)
			return status.Error(codes.Internal, "参数异常")
		}

		// 这里来构造update
		if dao2.NewInventoryDao(tx).StockIncrease(&incr) != nil {
			zap.L().Info(`<SellStock>.StockDecrease() != nil`)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (i InventoryServer) mustEmbedUnimplementedInventoryServer() {
	//TODO implement me
	panic("implement me")
}
