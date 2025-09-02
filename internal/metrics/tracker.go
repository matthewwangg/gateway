package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

type MetricsTracker struct {
	ServiceCalls *prometheus.CounterVec
	Requests     *prometheus.CounterVec
}

var Tracker *MetricsTracker

func Init() {
	Tracker = NewMetricsTracker()
}

func NewMetricsTracker() *MetricsTracker {
	tracker := &MetricsTracker{
		Requests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "requests",
			Help: "Number of requests to the gateway.",
		}, []string{"endpoint", "code"}),
		ServiceCalls: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "service_calls",
			Help: "Number of calls to each service.",
		}, []string{"service"}),
	}

	prometheus.MustRegister(
		tracker.Requests,
		tracker.ServiceCalls,
	)

	return tracker
}

func (t *MetricsTracker) RecordRequest(endpoint string, code int) {
	t.Requests.WithLabelValues(endpoint, strconv.Itoa(code)).Inc()
}

func (t *MetricsTracker) RecordServiceCall(service string) {
	t.ServiceCalls.WithLabelValues(service).Inc()
}
