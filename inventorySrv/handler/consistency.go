package handler

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"math/rand"
	dao2 "srv/inventorySrv/dao"
	"srv/inventorySrv/global"
	"srv/inventorySrv/model"
	"srv/inventorySrv/structs"
	"strconv"
	"time"
)

var (
	Pessimism = 0
	Optimism  = 1
	RedisLock = 2
)

type Consistency interface {
	//Select(...structs.Stocks) error
	Decr(...structs.Stocks) error // 减少库存
	Incr(...structs.Stocks) error // 归还库存
}

func GetConsistency(choice int) Consistency {
	switch choice {
	case 0:
		return NewPessimismLock()
	case 1:
		return NewOptimismLock()
	default:
		return NewRedLock()
	}
}

// 封装 1. 悲观锁  2. 乐观锁 3. redlock

// PessimismLock 悲观锁
type PessimismLock struct{}

func NewPessimismLock() *PessimismLock {
	return &PessimismLock{}
}

func (o *PessimismLock) Decr(info ...structs.Stocks) error {
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		// 用于快速IN查询
		goods := make([]int32, 0, len(info))
		// 用于快速查询是否足够
		goodsMap := make(map[int32]int32, len(info))
		// 用于调用封装的减库存
		decr := make([]dao2.Stock, 0, len(info))
		for _, v := range info {
			goods = append(goods, v.GoodsID)
			goodsMap[v.GoodsID] = v.Stock
			decr = append(decr, dao2.Stock{GoodsId: v.GoodsID, Stocks: v.Stock})
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

		for _, stock := range *infos {
			if stock.Stocks-goodsMap[stock.GoodsID] < 0 {
				return status.Error(codes.Internal, "商品库存不足")
			}
		}

		// 这里来构造update
		if dao2.NewInventoryDao(tx).StockDecrease(&decr) != nil {
			zap.L().Info(`<SellStock>.StockDecrease() != nil`)
			return status.Error(codes.Internal, "库存扣减失败")
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
func (o *PessimismLock) Incr(info ...structs.Stocks) error {
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		decr := make([]dao2.Stock, 0, len(info))
		for _, v := range info {
			decr = append(decr, dao2.Stock{GoodsId: v.GoodsID, Stocks: v.Stock})
		}

		// 这里来构造update
		if dao2.NewInventoryDao(tx).StockIncrease(&decr) != nil {
			zap.L().Info(`<SellStock>.StockDecrease() != nil`)
			return status.Error(codes.Internal, "库存扣减失败")
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// OptimismLock 乐观锁
type OptimismLock struct{}

func NewOptimismLock() *OptimismLock {
	return &OptimismLock{}
}

// Decr 扣减库存，如果失败就重试。没拿到就sleep rand.Intn(15)
func (o *OptimismLock) Decr(info ...structs.Stocks) error {
	// 1. 先进行数据的查询，额这里是不是还需要配置要重试的次数之类的东西
	// 2. 配置等待的时间
	// 3.
	if len(info) == 0 {
		return errors.New("参数异常")
	}
	tryTimes := 15
	if len(info) > tryTimes {
		tryTimes = len(info)
	}

	for i := 0; i < tryTimes; i++ {
		err := global.DB.Transaction(func(tx *gorm.DB) error {
			// 用于快速IN查询
			goods := make([]int32, 0, len(info))
			// 用于快速查询是否足够
			goodsMap := make(map[int32]int32, len(info))
			// 用于调用封装的减库存
			decr := make([]dao2.OptimismStock, 0, len(info))
			res := tx.Model(&model.Inventory{})
			for idx, v := range info {
				goods = append(goods, v.GoodsID)
				goodsMap[v.GoodsID] = v.Stock
				if idx == 0 {
					res.Where("goods_id = ? AND stocks >= ?", v.GoodsID, v.Stock)
					continue
				}
				res.Or("goods_id = ? AND stocks >= ?", v.GoodsID, v.Stock)
			}

			infos := &[]model.Inventory{}
			// 加读锁进行查询
			res = res.Find(&infos)
			if res.Error != nil {
				zap.L().Info("<SellStock>.Find(goods)", zap.Error(res.Error))
				return status.Error(codes.Internal, res.Error.Error())
			}

			if res.RowsAffected != int64(len(goods)) {
				zap.L().Info(`<SellStock>.RowsAffected != int64(len(goods))`)
				return status.Error(codes.Internal, "库存不足")
			}

			for _, stock := range *infos {
				// 这里特别要注意的是，是传递过来的stock
				decr = append(decr, dao2.OptimismStock{GoodsId: stock.GoodsID, Stocks: goodsMap[stock.GoodsID], Version: stock.Version})
				if stock.Stocks-goodsMap[stock.GoodsID] < 0 {
					return status.Error(codes.Internal, "库存不足")
				}
			}

			// 这里来构造update
			if dao2.NewInventoryDao(tx).OptimismDecr(&decr) != nil {
				zap.L().Info(`<SellStock>.StockDecrease() != nil`)
				return status.Error(codes.Internal, "库存扣减失败")
			}

			return nil
		})

		// 没有错误就返回
		if err == nil {
			break
		}

		if errors.Is(err, errors.New("库存不足")) {
			return status.Error(codes.InvalidArgument, "库存不足")
		}

		time.Sleep(time.Millisecond * time.Duration(rand.Intn(15)))
	}

	return nil
}
func (o *OptimismLock) Incr(info ...structs.Stocks) error {
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		decr := make([]dao2.Stock, 0, len(info))
		for _, v := range info {
			decr = append(decr, dao2.Stock{GoodsId: v.GoodsID, Stocks: v.Stock})
		}

		// 这里来构造update
		if dao2.NewInventoryDao(tx).StockIncrease(&decr) != nil {
			zap.L().Info(`<SellStock>.StockDecrease() != nil`)
			return status.Error(codes.Internal, "库存扣减失败")
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// RedLock 红锁
type RedLock struct{}

func NewRedLock() *RedLock {
	return &RedLock{}
}

func (o *RedLock) Decr(info ...structs.Stocks) error {
	// 1. 应该是查询的时候先获取到锁，然后查询完成之后释放锁，中途不需要加乐观锁或者其他什么东西的
	// Obtain a new mutex by using the same name for all instances wanting the
	// same lock.
	if len(info) == 0 {
		return status.Error(codes.InvalidArgument, "参数异常")
	}
	mutexname := "my-global-mutex"
	for i, stocks := range info {
		mutexname += strconv.Itoa(int(stocks.GoodsID))
		if i < len(info)-1 {
			mutexname += ","
		}
	}

	mutex := global.RedLock.NewMutex(mutexname)

	// Obtain a lock for our given mutex. After this is successful, no one else
	// can obtain the same lock (the same mutex name) until we unlock it.
	if err := mutex.Lock(); err != nil {
		return status.Error(codes.Internal, fmt.Sprintf("RedLock.Lock Get Failed: %s", err.Error()))
	}

	err := global.DB.Transaction(func(tx *gorm.DB) error {
		// 用于快速IN查询
		goods := make([]int32, 0, len(info))
		// 用于快速查询是否足够
		goodsMap := make(map[int32]int32, len(info))
		// 用于调用封装的减库存
		decr := make([]dao2.OptimismStock, 0, len(info))
		res := tx.Model(&model.Inventory{})
		for idx, v := range info {
			goods = append(goods, v.GoodsID)
			goodsMap[v.GoodsID] = v.Stock
			if idx == 0 {
				res.Where("goods_id = ? AND stocks >= ?", v.GoodsID, v.Stock)
				continue
			}
			res.Or("goods_id = ? AND stocks >= ?", v.GoodsID, v.Stock)
		}

		infos := &[]model.Inventory{}
		res = res.Find(&infos)
		if res.Error != nil {
			zap.L().Info("(o *RedLock) Decr(info ...structs.Stocks) error", zap.Error(res.Error))
			return status.Error(codes.Internal, res.Error.Error())
		}

		if res.RowsAffected != int64(len(goods)) {
			zap.L().Info(`<(o *RedLock) Decr(info ...structs.Stocks) error>.RowsAffected != int64(len(goods))`)
			return status.Error(codes.Internal, "库存不足")
		}

		for _, stock := range *infos {
			// 这里特别要注意的是，是传递过来的stock
			decr = append(decr, dao2.OptimismStock{GoodsId: stock.GoodsID, Stocks: goodsMap[stock.GoodsID]})
			if stock.Stocks-goodsMap[stock.GoodsID] < 0 {
				return status.Error(codes.Internal, "库存不足")
			}
		}

		// 这里来构造update
		if dao2.NewInventoryDao(tx).LockDecr(&decr) != nil {
			zap.L().Info(`<SellStock>.StockDecrease() != nil`)
			return status.Error(codes.Internal, "库存扣减失败")
		}

		// Do your work that requires the lock.

		// Release the lock so other processes or threads can obtain a lock.
		if ok, err := mutex.Unlock(); !ok || err != nil {
			return status.Error(codes.Internal, fmt.Sprintf("RedLock.Lock Get Failed: %s", err.Error()))
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
func (o *RedLock) Incr(info ...structs.Stocks) error {
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		decr := make([]dao2.Stock, 0, len(info))
		for _, v := range info {
			decr = append(decr, dao2.Stock{GoodsId: v.GoodsID, Stocks: v.Stock})
		}

		// 这里来构造update
		if dao2.NewInventoryDao(tx).StockIncrease(&decr) != nil {
			zap.L().Info(`<SellStock>.StockDecrease() != nil`)
			return status.Error(codes.Internal, "库存扣减失败")
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
