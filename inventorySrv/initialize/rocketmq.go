package initialize

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"log"
	"os"
	dao2 "srv/inventorySrv/dao"
	"srv/inventorySrv/global"
	"srv/inventorySrv/handler"
	"srv/inventorySrv/structs"
	"time"

	"github.com/apache/rocketmq-clients/golang"
	"github.com/apache/rocketmq-clients/golang/credentials"
)

const (
	Endpoint      = "192.168.3.5:8081"
	AccessKey     = "xxxxxx"
	SecretKey     = "xxxxxx"
	TrxTopic      = "trx_msg_rollback_stock"
	ConsumerGroup = "trx_msg_rollback_stock"
)

var (
	// maximum waiting time for receive func
	awaitDuration = time.Second * 5
	// maximum number of messages received at one time
	maxMessageNum int32 = 16
	// invisibleDuration should > 20s
	invisibleDuration = time.Second * 20
	// receive messages in a loop
)

func InitSubscribe() {
	zap.L().Info("START InitSubscribe")
	// log to console
	os.Setenv("mq.consoleAppender.enabled", "true")
	golang.ResetLogger()
	// new simpleConsumer instance
	simpleConsumer, err := golang.NewSimpleConsumer(&golang.Config{
		Endpoint:      Endpoint,
		ConsumerGroup: ConsumerGroup,
		Credentials: &credentials.SessionCredentials{
			AccessKey:    AccessKey,
			AccessSecret: SecretKey,
		},
	},
		golang.WithAwaitDuration(awaitDuration),
		golang.WithSubscriptionExpressions(map[string]*golang.FilterExpression{
			TrxTopic: golang.SUB_ALL,
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	// start simpleConsumer
	err = simpleConsumer.Start()
	if err != nil {
		log.Fatal(err)
	}
	// graceful stop simpleConsumer
	defer simpleConsumer.GracefulStop()
	zap.L().Info("START InitSubscribe ok")
	for {
		fmt.Println("start recevie message")
		mvs, err := simpleConsumer.Receive(context.TODO(), maxMessageNum, invisibleDuration)
		if err != nil {
			fmt.Println(err)
		}
		// ack message
		zap.L().Info("read from :", zap.Any("mv", mvs))
		for _, mv := range mvs {
			/*
				调用下面的方法来实现库存的归还。以及归还完成之后，对订单归还状态进行修改。
				当切仅当错误是数据库内部错误，或者是归还成功才进行 ACK 应答
			*/
			orderSN := string(mv.GetBody())
			// 调用归还库存的模块
			resp, err := dao2.NewOrderDetailDao(global.DB).GetOne(orderSN)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}

			if errors.Is(err, gorm.ErrRecordNotFound) {
				simpleConsumer.Ack(context.TODO(), mv)
				continue
			}

			stock := make([]structs.Stocks, 0, len(resp.Details))
			for _, detail := range resp.Details {
				stock = append(stock, structs.Stocks{
					GoodsID: detail.GoodsID,
					Stock:   detail.GoodsNum,
					Version: 0,
				})
			}

			// 调用归还库存的接口
			err = handler.GetConsistency(999).Incr(stock...)
			if err != nil {
				zap.L().Error("事务自动归还库存失败", zap.Error(err))
				continue
			}

			_ = dao2.NewOrderDetailDao(global.DB).UpdateStatus(orderSN)

			// 如果是找到了，则选择归还，成功则提交，否则不ack
			simpleConsumer.Ack(context.TODO(), mv)
			fmt.Println(mv)
		}
		fmt.Println("wait a moment")
		fmt.Println()
		time.Sleep(time.Second * 1)
	}
}
