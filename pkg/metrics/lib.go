package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var SmsStatusMetric *prometheus.CounterVec

var registry = prometheus.NewRegistry()

func RegisterMetrics() {
	SmsStatusMetric = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "message_status_count",
		Help: "The total number of processed messages",
	}, []string{"status"})

	registry.MustRegister(SmsStatusMetric)
}

func GetRegistry() *prometheus.Registry {
	return registry
}
