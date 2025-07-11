package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	RequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_requests_total",
			Help: "total number of gRPC request for service",
		},
		[]string{"crud_type"},
	)

	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "grpc_request_duration",
			Help: "Duration for route request in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"crud_type"},
	)

	StartTime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "start_time",
			Help: "Time to start app",
			Buckets: prometheus.DefBuckets,
		},
		[]string{},
	)
)

func init() {
	prometheus.MustRegister(RequestCount)
	prometheus.MustRegister(RequestDuration)
	prometheus.MustRegister(StartTime)
}