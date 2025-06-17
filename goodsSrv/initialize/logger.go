package initialize

import (
	"go.uber.org/zap"
)

func InitZap() {
	dev, err := zap.NewDevelopment()
	if err != nil {
		panic("init zap err:" + err.Error())
	}
	zap.ReplaceGlobals(dev)
}
