package bootstrap

import (
	"context"
	"errors"

	conf "github.com/fzf-labs/kratos-contrib/api/conf/v1"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	traceSdk "go.opentelemetry.io/otel/sdk/trace"
	semConv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// NewTracerProvider 创建一个链路追踪器
func NewTracerProvider(cfg *conf.Tracer, serviceInfo *Service) error {
	if cfg == nil {
		return errors.New("tracer config is nil")
	}
	if cfg.Sampler == 0 {
		cfg.Sampler = 1.0
	}
	opts := []traceSdk.TracerProviderOption{
		// 将基于父span的采样率设置为Sampler(1.0=100%)
		traceSdk.WithSampler(traceSdk.ParentBased(traceSdk.TraceIDRatioBased(cfg.Sampler))),
		traceSdk.WithResource(resource.NewSchemaless(
			semConv.ServiceNameKey.String(serviceInfo.Name),
			semConv.ServiceVersionKey.String(serviceInfo.Version),
			semConv.ServiceInstanceIDKey.String(serviceInfo.ID),
		)),
	}
	if len(cfg.Endpoint) > 0 {
		// 初始化采集器
		exp, err := NewTracerExporter(cfg.Batcher, cfg.Endpoint, cfg.Insecure)
		if err != nil {
			panic(err)
		}
		// 始终确保在生产中批量处理
		opts = append(opts, traceSdk.WithBatcher(exp))
	}
	tp := traceSdk.NewTracerProvider(opts...)
	if tp == nil {
		return errors.New("create tracer provider failed")
	}
	otel.SetTracerProvider(tp)
	return nil
}

// NewTracerExporter 创建一个导出器，支持：zipkin、otlp-http、otlp-grpc
func NewTracerExporter(exporterName, endpoint string, insecure bool) (traceSdk.SpanExporter, error) {
	ctx := context.Background()
	switch exporterName {
	case "otlphttp":
		return NewOtlpHTTPExporter(ctx, endpoint, insecure)
	case "otlpgrpc":
		return NewOtlpGrpcExporter(ctx, endpoint, insecure)
	case "stdout":
		return NewStdoutExporter()
	default:
		return NewStdoutExporter()
	}
}

// NewStdoutExporter 创建一个标准输出导出器
func NewStdoutExporter() (traceSdk.SpanExporter, error) {
	return stdouttrace.New()
}

// NewOtlpHTTPExporter 创建OTLP/HTTP导出器，默认端口：4318
func NewOtlpHTTPExporter(ctx context.Context, endpoint string, insecure bool, options ...otlptracehttp.Option) (traceSdk.SpanExporter, error) {
	var opts []otlptracehttp.Option
	opts = append(opts, otlptracehttp.WithEndpoint(endpoint))

	if insecure {
		opts = append(opts, otlptracehttp.WithInsecure())
	}

	opts = append(opts, options...)

	return otlptrace.New(
		ctx,
		otlptracehttp.NewClient(opts...),
	)
}

// NewOtlpGrpcExporter 创建OTLP/gRPC导出器，默认端口：4317
func NewOtlpGrpcExporter(ctx context.Context, endpoint string, insecure bool, options ...otlptracegrpc.Option) (traceSdk.SpanExporter, error) {
	var opts []otlptracegrpc.Option
	opts = append(opts, otlptracegrpc.WithEndpoint(endpoint))

	if insecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}

	opts = append(opts, options...)

	return otlptrace.New(
		ctx,
		otlptracegrpc.NewClient(opts...),
	)
}
