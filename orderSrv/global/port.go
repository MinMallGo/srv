package global

import (
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"math/rand"
	"net"
	"time"
)

// GetPort 获取能用的随机端口
func GetPort() int {
	// 通过 :0 来获取随机端口
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		zap.L().Fatal("net.Listen", zap.Error(err))
	}
	defer listener.Close()

	addr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		zap.L().Fatal("assert *net.TCPAddr with error.", zap.Error(err))
	}

	return addr.Port
}

func UUID() string {
	return uuid.New().String()
}

func OrderSN(userId int) string {
	now := time.Now()
	randx := rand.New(rand.NewSource(now.UnixNano())).Intn(90) + 10

	return fmt.Sprintf("%s%d%d%d", now.Format("20060102150405"), now.UnixNano(), userId, randx)
}
