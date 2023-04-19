package comet

import "github.com/go-kratos/kratos/pkg/stat/metric"

const (
	serverNamespace = "comet_server"
)

var (
	_metricServerReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "comet server requests duration(ms).",
		Labels:    []string{"method", "caller"},
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	})
	_metricServerReqCodeTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "comet server requests code count.",
		Labels:    []string{"method", "caller", "code"},
	})
)

