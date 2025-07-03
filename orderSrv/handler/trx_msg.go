package handler

import (
	"context"
	"errors"
	rmq_client "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/apache/rocketmq-clients/golang/v5/credentials"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	"srv/orderSrv/dao"
	"srv/orderSrv/global"
)

const (
	Endpoint   = "192.168.3.5:8081"
	AccessKey  = "xxxxxx"
	SecretKey  = "xxxxxx"
	TrxTopic   = "trx_msg_rollback_stock"
	DelayTopic = "xxxxxx"
)

func checker(msg *rmq_client.MessageView) rmq_client.TransactionResolution {
	// 这里是你的业务逻辑，用于查询本地事务的最终状态
	// 通常会根据 msg.TransactionId 或 msg.Keys 去查询数据库或其他持久化存储
	// 例如：检查订单状态是否已成功创建
	//
	// 为了测试，我们模拟不同的回查结果：

	// 1. 模拟本地事务确实成功了，通知COMMIT
	// return primitive.TransactionCommit

	// 2. 模拟本地事务确实失败了，通知ROLLBACK
	// return primitive.TransactionRollback

	// 3. 模拟本地事务状态不确定，Broker会稍后再次回查
	// 	return rmq_client.UNKNOWN

	// 这里根据消息队列提供的orderSN来查询选择返回哪个信息
	orderSN := string(msg.GetBody())
	zap.L().Info("<正在触发补偿机制>", zap.String("orderSN", orderSN))
	// 这里就查询本地的订单是否存在就行了，因为本地订单是最后才创建的，如果创建了，说明前面的流程全部都走完了
	detail, err := dao.NewOrderDao(global.DB).GetDetail(context.Background(), dao.OrderDetailReq{OrderSN: orderSN})
	if err != nil && errors.Is(err, dao.DBErr) {
		return rmq_client.UNKNOWN
	}

	// TODO 这里应该是提交？ 如果本地没有查询到订单的话，就需要归还库存
	// 因为不清楚库存服务是否扣减了库存，所以要归还，然后这里应该提交这条事务消息，因为对面好像要订阅
	if detail == nil || detail.ID == 0 {
		return rmq_client.COMMIT
	}

	return rmq_client.ROLLBACK
}

func NewTrxMsg() (rmq_client.Producer, error) {
	os.Setenv("mq.consoleAppender.enabled", "true")
	rmq_client.ResetLogger()
	producer, err := rmq_client.NewProducer(
		&rmq_client.Config{
			Endpoint: Endpoint,
			Credentials: &credentials.SessionCredentials{
				AccessKey:    AccessKey,
				AccessSecret: SecretKey,
			},
		},
		rmq_client.WithTopics(TrxTopic), // 指定topic
		rmq_client.WithTransactionChecker(&rmq_client.TransactionChecker{ // 带上事务
			Check: checker,
		}),
	)

	if err != nil {
		zap.L().Error("初始化事务消息失败", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	err = producer.Start()
	if err != nil {
		zap.L().Error("启动事务消息失败", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return producer, err
}
