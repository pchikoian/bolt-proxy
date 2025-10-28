package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Connection metrics
	ActiveConnections = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "bolt_proxy_active_connections",
		Help: "Number of currently active client connections",
	})

	TotalConnections = promauto.NewCounter(prometheus.CounterOpts{
		Name: "bolt_proxy_connections_total",
		Help: "Total number of client connections",
	})

	ConnectionDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "bolt_proxy_connection_duration_seconds",
		Help:    "Duration of client connections in seconds",
		Buckets: prometheus.DefBuckets,
	})

	// Authentication metrics
	AuthAttempts = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "bolt_proxy_auth_attempts_total",
		Help: "Total number of authentication attempts",
	}, []string{"status"}) // status: success, failure

	AuthDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "bolt_proxy_auth_duration_seconds",
		Help:    "Duration of authentication attempts in seconds",
		Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1},
	})

	// Message metrics
	MessagesForwarded = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "bolt_proxy_messages_forwarded_total",
		Help: "Total number of messages forwarded",
	}, []string{"direction", "type"}) // direction: client_to_server, server_to_client; type: HELLO, RUN, etc.

	MessageBytes = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "bolt_proxy_message_bytes_total",
		Help: "Total bytes of messages forwarded",
	}, []string{"direction"})

	MessageProcessingDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "bolt_proxy_message_processing_duration_seconds",
		Help:    "Duration of message processing in seconds",
		Buckets: []float64{.0001, .0005, .001, .005, .01, .025, .05, .1},
	}, []string{"type"})

	// Backend metrics
	BackendConnections = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "bolt_proxy_backend_connections",
		Help: "Number of active connections to backend servers",
	}, []string{"backend"})

	BackendConnectionErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "bolt_proxy_backend_connection_errors_total",
		Help: "Total number of backend connection errors",
	}, []string{"backend", "error_type"})

	BackendLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "bolt_proxy_backend_latency_seconds",
		Help:    "Latency of backend connections in seconds",
		Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5},
	}, []string{"backend"})

	// Transaction metrics
	ActiveTransactions = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "bolt_proxy_active_transactions",
		Help: "Number of currently active transactions",
	})

	TransactionsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "bolt_proxy_transactions_total",
		Help: "Total number of transactions",
	}, []string{"status"}) // status: committed, rolled_back, failed

	// Error metrics
	Errors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "bolt_proxy_errors_total",
		Help: "Total number of errors by type",
	}, []string{"error_type", "component"})

	// Health check metrics
	HealthChecks = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "bolt_proxy_health_checks_total",
		Help: "Total number of health check requests",
	}, []string{"status"})

	// Protocol metrics
	BoltVersionNegotiations = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "bolt_proxy_bolt_version_negotiations_total",
		Help: "Total number of Bolt version negotiations",
	}, []string{"version"})

	// CRUD operation metrics
	CRUDOperations = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "bolt_proxy_crud_operations_total",
		Help: "Total number of CRUD operations by type",
	}, []string{"operation"}) // operation: create, read, update, delete, other

	// Query metrics (Requests Per Query - RPQ)
	QueryExecutions = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "bolt_proxy_query_executions_total",
		Help: "Total number of query executions",
	}, []string{"type"}) // type: run, pull, discard

	QueryDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "bolt_proxy_query_duration_seconds",
		Help:    "Duration of query execution in seconds",
		Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
	}, []string{"type"})
)

// Helper functions to simplify metric recording

func RecordConnection() {
	TotalConnections.Inc()
	ActiveConnections.Inc()
}

func RecordConnectionClosed() {
	ActiveConnections.Dec()
}

func RecordAuthSuccess() {
	AuthAttempts.WithLabelValues("success").Inc()
}

func RecordAuthFailure() {
	AuthAttempts.WithLabelValues("failure").Inc()
}

func RecordMessageForwarded(direction, msgType string, bytes int) {
	MessagesForwarded.WithLabelValues(direction, msgType).Inc()
	MessageBytes.WithLabelValues(direction).Add(float64(bytes))
}

func RecordError(errorType, component string) {
	Errors.WithLabelValues(errorType, component).Inc()
}

func RecordHealthCheck(status string) {
	HealthChecks.WithLabelValues(status).Inc()
}

func RecordBoltVersion(version string) {
	BoltVersionNegotiations.WithLabelValues(version).Inc()
}

func RecordBackendConnection(backend string) {
	BackendConnections.WithLabelValues(backend).Inc()
}

func RecordBackendConnectionClosed(backend string) {
	BackendConnections.WithLabelValues(backend).Dec()
}

func RecordBackendConnectionError(backend, errorType string) {
	BackendConnectionErrors.WithLabelValues(backend, errorType).Inc()
}

func RecordTransactionStart() {
	ActiveTransactions.Inc()
}

func RecordTransactionEnd(status string) {
	ActiveTransactions.Dec()
	TransactionsTotal.WithLabelValues(status).Inc()
}

func RecordCRUDOperation(operation string) {
	CRUDOperations.WithLabelValues(operation).Inc()
}

func RecordQueryExecution(queryType string) {
	QueryExecutions.WithLabelValues(queryType).Inc()
}
