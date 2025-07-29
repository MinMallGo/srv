package utils

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	traceSDK "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"time"
)

const (
	serviceName    = "Go-Jaeger-Demo"
	jaegerEndpoint = "127.0.0.1:4318" // todo 这里需要配置到nacos
)

// SetupTracer 设置Tracer
func SetupTracer(ctx context.Context) (func(context.Context) error, error) {
	tracerProvider, err := newJaegerTraceProvider(ctx)
	if err != nil {
		return nil, err
	}
	otel.SetTracerProvider(tracerProvider)
	return tracerProvider.Shutdown, nil
}

// NewJaegerTraceProvider 创建一个 Jaeger Trace Provider
func newJaegerTraceProvider(ctx context.Context) (*traceSDK.TracerProvider, error) {
	// 创建一个使用 HTTP 协议连接本机Jaeger的 Exporter
	exp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(jaegerEndpoint),
		otlptracehttp.WithInsecure())
	if err != nil {
		return nil, err
	}
	res, err := resource.New(ctx, resource.WithAttributes(semconv.ServiceName(serviceName)))
	if err != nil {
		return nil, err
	}
	traceProvider := traceSDK.NewTracerProvider(
		traceSDK.WithResource(res),
		traceSDK.WithSampler(traceSDK.AlwaysSample()), // 采样
		traceSDK.WithBatcher(exp, traceSDK.WithBatchTimeout(time.Second)),
	)
	return traceProvider, nil
}

func WithSpan[T any](ctx context.Context, tracer trace.Tracer, name string, attrs []attribute.KeyValue, fn func(context.Context) (T, error)) (T, error) {
	var span trace.Span
	ctx, span = tracer.Start(ctx, name)
	defer span.End()

	res, err := fn(ctx)

	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return res, err
	}

	return res, nil
}
