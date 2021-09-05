package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics interface {
	IncrementCreateCounter()
	IncrementMultiCreateCounter()
	IncrementUpdateCounter()
	IncrementRemoveCounter()
}

type PrometheusMetrics struct {
	createCounter      prometheus.Counter
	multiCreateCounter prometheus.Counter
	updateCounter      prometheus.Counter
	removeCounter      prometheus.Counter
}

func NewPrometheusMetrics() *PrometheusMetrics {
	createCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "grpc_request_create_succeed_count",
		Help: "Number of successfully handled Create requests",
	})

	multiCreateCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "grpc_request_multi_create_succeed_count",
		Help: "Number of successfully handled MultiCreate requests",
	})

	updateCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "grpc_request_update_succeed_count",
		Help: "Number of successfully handled Update requests",
	})

	removeCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "grpc_request_remove_succeed_count",
		Help: "Number of successfully handled Remove requests",
	})

	prometheus.MustRegister(createCounter, multiCreateCounter, updateCounter, removeCounter)

	return &PrometheusMetrics{
		createCounter:      createCounter,
		multiCreateCounter: multiCreateCounter,
		updateCounter:      updateCounter,
		removeCounter:      removeCounter,
	}
}

func (m *PrometheusMetrics) IncrementCreateCounter() {
	m.createCounter.Inc()
}

func (m *PrometheusMetrics) IncrementMultiCreateCounter() {
	m.multiCreateCounter.Inc()
}

func (m *PrometheusMetrics) IncrementUpdateCounter() {
	m.updateCounter.Inc()
}

func (m *PrometheusMetrics) IncrementRemoveCounter() {
	m.removeCounter.Inc()
}
