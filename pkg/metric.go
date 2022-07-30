package main

import "github.com/prometheus/client_golang/prometheus"

var (
	handlingSecondsBucket30 = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "grpc",
		Subsystem: "server",
		Help:      "num of alert task error",
		Name:      "handling_seconds_bucket_30s",
	})

	handlingSecondsBucket50 = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "grpc",
		Subsystem: "server",
		Help:      "num of alert task error",
		Name:      "handling_seconds_bucket_50s",
	})

	handlingSecondsBucket60 = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "grpc",
		Subsystem: "server",
		Help:      "num of alert task error",
		Name:      "handling_seconds_bucket_60s",
	})

	queryPrestoCost = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "presto",
		Subsystem: "query",
		Help:      "query presto cost millisecond",
		Name:      "millisecond",
		Buckets:   []float64{1000, 5000, 10000, 30000, 60000, 120000},
	})

	tokenUpdateError = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "presto",
		Subsystem: "plugin",
		Help:      "num of update token from token service error",
		Name:      "token_update_error",
	}, []string{"token_server"})

	queryError = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "presto",
		Subsystem: "plugin",
		Help:      "num of update query error",
		Name:      "query_error",
	})
)

func init() {
	prometheus.MustRegister(
		handlingSecondsBucket30,
		handlingSecondsBucket50,
		handlingSecondsBucket60,
		queryPrestoCost,
		tokenUpdateError,
		queryError,
	)
}
