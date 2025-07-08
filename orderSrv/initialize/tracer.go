package initialize

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"srv/orderSrv/global"
	"srv/orderSrv/utils"
)

func InitTracing() {
	ctx := context.Background()
	shutdown, err := utils.SetupTracer(ctx)
	if err != nil {
		panic(fmt.Errorf("tracer setup failed: %v", err))
	}
	otel.SetTextMapPropagator(propagation.TraceContext{})
	// 用于优雅退出
	global.TraceShutdownFunc = shutdown
}
