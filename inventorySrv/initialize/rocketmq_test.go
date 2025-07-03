package initialize

import (
	"testing"
)

func TestMQTrx(t *testing.T) {
	initxxxx()
	InitSubscribe()

	// 往这里发送几条消息
}

func initxxxx() {
	InitZap()
	InitConfig() // 如果要测试，这里要改一下
	InitDB()
	InitRedLock() // 初始化分布式锁
}
