package utils

import (
	"context"
	"github.com/apache/rocketmq-clients/golang"
	rmq_client "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/apache/rocketmq-clients/golang/v5/credentials"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	"time"
)

const (
	Endpoint    = "192.168.3.5:8081"
	AccessKey   = "xxxxxx"
	SecretKey   = "xxxxxx"
	TrxTopic    = "trx_msg_rollback_stock"
	DelayTopic  = "xxxxxx"
	ExpireTopic = "order_delay_cancel"
)

type RocketMQ struct{}

type MQResult struct {
	Product rmq_client.Producer
	Trx     golang.Transaction
	Receipt []*rmq_client.SendReceipt
}

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

	// TODO 根据业务逻辑来检查是否提交该条消息。判断订单是否创建成功。创建成功则commit，没检查到就unknown

	return rmq_client.ROLLBACK
}

func NewTrxMsg() (MQResult, error) {
	os.Setenv("mq.consoleAppender.enabled", "true")
	rmq_client.ResetLogger()
	producer, err := rmq_client.NewProducer(&rmq_client.Config{
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

	mqResult := MQResult{
		Product: producer,
	}

	if err != nil {
		zap.L().Error("初始化事务消息失败", zap.Error(err))
		return mqResult, status.Error(codes.Internal, err.Error())
	}
	err = producer.Start()
	if err != nil {
		zap.L().Error("启动事务消息失败", zap.Error(err))
		return mqResult, status.Error(codes.Internal, err.Error())
	}

	return mqResult, err
}

func SendTrxMsg(ctx context.Context, product rmq_client.Producer, topic, msg string) (MQResult, error) {
	trx := product.BeginTransaction()
	message := &rmq_client.Message{
		Topic: topic,
		Body:  []byte(msg),
	}
	message.SetTag("order_generate")
	resp, err := product.SendWithTransaction(ctx, message, trx)
	if err != nil {
		return MQResult{}, err
	}
	return MQResult{Trx: trx, Product: product, Receipt: resp}, nil
}

func SendOrderDelayMsg(orderSN string) error {
	// log to console
	os.Setenv("mq.consoleAppender.enabled", "true")
	rmq_client.ResetLogger()
	// In most case, you don't need to create many producers, singleton pattern is more recommended.
	producer, err := rmq_client.NewProducer(&rmq_client.Config{
		Endpoint: Endpoint,
		Credentials: &credentials.SessionCredentials{
			AccessKey:    AccessKey,
			AccessSecret: SecretKey,
		},
	},
		rmq_client.WithTopics(ExpireTopic),
	)
	if err != nil {
		zap.L().Error("初始化延时消息失败", zap.Error(err))
		return status.Error(codes.Internal, err.Error())
	}
	// start producer
	err = producer.Start()
	if err != nil {
		zap.L().Error("启动延时消息失败", zap.Error(err))
		return status.Error(codes.Internal, err.Error())
	}
	// graceful stop producer
	defer producer.GracefulStop()

	// new a message
	msg := &rmq_client.Message{
		Topic: ExpireTopic,
		Body:  []byte(orderSN),
	}
	// set keys and tag
	msg.SetKeys("order_sn", orderSN)
	msg.SetTag("order")
	// set delay timestamp
	msg.SetDelayTimestamp(time.Now().Add(time.Second * 30))
	// send message in sync
	resp, err := producer.Send(context.TODO(), msg)
	if err != nil {
		zap.L().Error("发送订单延时取消消息失败：", zap.Error(err))
	}

	zap.L().Info("发送订单延时取消消息成功", zap.Any("resp", resp))

	return nil
}
