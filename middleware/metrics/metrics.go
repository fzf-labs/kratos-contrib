package metrics

import (
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// Detailed reference https://github.com/go-kratos/examples/tree/main/metrics
var (
	_metricRequests metric.Int64Counter
	_metricSeconds  metric.Float64Histogram
)

func init() {
	exporter, err := prometheus.New()
	if err != nil {
		panic(err)
	}
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(exporter))
	meter := provider.Meter("metrics")
	_metricRequests, err = metrics.DefaultRequestsCounter(meter, metrics.DefaultServerRequestsCounterName)
	if err != nil {
		panic(err)
	}
	_metricSeconds, err = metrics.DefaultSecondsHistogram(meter, metrics.DefaultServerSecondsHistogramName)
	if err != nil {
		panic(err)
	}
}

func Server() middleware.Middleware {
	return metrics.Server(
		metrics.WithRequests(_metricRequests),
		metrics.WithSeconds(_metricSeconds),
	)
}
