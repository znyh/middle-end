package mongdb

import (
	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "mongodb_client"

var (
	_metricReqDur = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "mongodb client requests duration(ms).",
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500},
	}, []string{"name", "addr", "command"})

	_metricReqErr = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "requests",
		Name:      "error_total",
		Help:      "mongodb client requests error count.",
	}, []string{"name", "addr", "command", "error"})
	_metricConnTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "connections",
		Name:      "total",
		Help:      "mongodb client connections total count.",
	}, []string{"name", "addr", "state"})
	_metricConnCurrent = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "connections",
		Name:      "current",
		Help:      "mongodb client connections current.",
	}, []string{"name", "addr", "state"})
)
