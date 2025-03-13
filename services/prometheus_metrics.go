package services

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// PrometheusMetricsInterface definieert de interface voor metrics
type PrometheusMetricsInterface interface {
	RecordEmailSent(emailType, status string)
	RecordEmailFailed(emailType, reason string)
	ObserveEmailLatency(emailType string, duration float64)
}

// PrometheusMetrics verzorgt Prometheus monitoring
type PrometheusMetrics struct {
	emailsSent           *prometheus.CounterVec
	emailsFailed         *prometheus.CounterVec
	emailLatency         *prometheus.HistogramVec
	rateLimitExceeded    *prometheus.CounterVec
	activeEmailBatches   prometheus.Gauge
	mu                   sync.Mutex
	emailTypeCardinality map[string]bool // Helpt bij het beperken van cardinality
}

var (
	prometheusMetricsInstance *PrometheusMetrics
	prometheusMetricsMutex    sync.Mutex
)

// GetPrometheusMetrics geeft een singleton instantie van PrometheusMetrics terug
func GetPrometheusMetrics() *PrometheusMetrics {
	prometheusMetricsMutex.Lock()
	defer prometheusMetricsMutex.Unlock()

	if prometheusMetricsInstance == nil {
		prometheusMetricsInstance = newPrometheusMetrics()
	}

	return prometheusMetricsInstance
}

// newPrometheusMetrics is nu private
func newPrometheusMetrics() *PrometheusMetrics {
	// Definieer email verzending metrics
	emailsSent := promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "email_service_emails_sent_total",
		Help: "Het totale aantal succesvol verzonden emails",
	}, []string{"type", "template"})

	emailsFailed := promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "email_service_emails_failed_total",
		Help: "Het totale aantal mislukte emails",
	}, []string{"type", "error_type"})

	// Email verwerkingstijd (latency)
	emailLatency := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "email_service_latency_seconds",
		Help: "Tijd benodigd voor het verzenden van een email",
		// Buckets specifiek voor email verzending (in seconden)
		Buckets: []float64{0.1, 0.5, 1.0, 2.5, 5.0, 10.0},
	}, []string{"type"})

	// Rate limiting metrics
	rateLimitExceeded := promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "email_service_rate_limit_exceeded_total",
		Help: "Het aantal keer dat een rate limit is overschreden",
	}, []string{"type", "limit_type"})

	// Actieve email batches
	activeEmailBatches := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "email_service_active_batches",
		Help: "Het huidige aantal actieve email batches",
	})

	return &PrometheusMetrics{
		emailsSent:           emailsSent,
		emailsFailed:         emailsFailed,
		emailLatency:         emailLatency,
		rateLimitExceeded:    rateLimitExceeded,
		activeEmailBatches:   activeEmailBatches,
		emailTypeCardinality: make(map[string]bool),
	}
}

// RecordEmailSent registreert een succesvol verzonden email in Prometheus
func (pm *PrometheusMetrics) RecordEmailSent(emailType, template string) {
	pm.checkEmailTypeCardinality(emailType)
	pm.emailsSent.WithLabelValues(emailType, template).Inc()
}

// RecordEmailFailed registreert een mislukte email in Prometheus
func (pm *PrometheusMetrics) RecordEmailFailed(emailType, errorType string) {
	pm.checkEmailTypeCardinality(emailType)
	pm.emailsFailed.WithLabelValues(emailType, errorType).Inc()
}

// ObserveEmailLatency registreert de duur van een email verzending
func (pm *PrometheusMetrics) ObserveEmailLatency(emailType string, duration float64) {
	pm.checkEmailTypeCardinality(emailType)
	pm.emailLatency.WithLabelValues(emailType).Observe(duration)
}

// RecordRateLimitExceeded registreert een rate limit overschrijding
func (pm *PrometheusMetrics) RecordRateLimitExceeded(emailType, limitType string) {
	pm.checkEmailTypeCardinality(emailType)
	pm.rateLimitExceeded.WithLabelValues(emailType, limitType).Inc()
}

// UpdateActiveBatches werkt het aantal actieve batches bij
func (pm *PrometheusMetrics) UpdateActiveBatches(count int) {
	pm.activeEmailBatches.Set(float64(count))
}

// checkEmailTypeCardinality voorkomt cardinality explosie door onbekende emailTypes
func (pm *PrometheusMetrics) checkEmailTypeCardinality(emailType string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Als we dit email type nog niet eerder hebben gezien, log het
	if _, exists := pm.emailTypeCardinality[emailType]; !exists {
		// In een echte implementatie zou je hier een warning kunnen loggen
		// of een maximum aantal types kunnen handhaven
		pm.emailTypeCardinality[emailType] = true
	}
}

// NewPrometheusMetricsWithRegistry maakt een nieuwe Prometheus metrics instance met een aangepaste registry
func NewPrometheusMetricsWithRegistry(reg *prometheus.Registry) *PrometheusMetrics {
	factory := promauto.With(reg)

	// Definieer email verzending metrics
	emailsSent := factory.NewCounterVec(prometheus.CounterOpts{
		Name: "email_service_emails_sent_total",
		Help: "Het totale aantal succesvol verzonden emails",
	}, []string{"type", "template"})

	emailsFailed := factory.NewCounterVec(prometheus.CounterOpts{
		Name: "email_service_emails_failed_total",
		Help: "Het totale aantal mislukte emails",
	}, []string{"type", "error_type"})

	// Email verwerkingstijd (latency)
	emailLatency := factory.NewHistogramVec(prometheus.HistogramOpts{
		Name: "email_service_latency_seconds",
		Help: "Tijd benodigd voor het verzenden van een email",
		// Buckets specifiek voor email verzending (in seconden)
		Buckets: []float64{0.1, 0.5, 1.0, 2.5, 5.0, 10.0},
	}, []string{"type"})

	// Rate limiting metrics
	rateLimitExceeded := factory.NewCounterVec(prometheus.CounterOpts{
		Name: "email_service_rate_limit_exceeded_total",
		Help: "Het aantal keer dat een rate limit is overschreden",
	}, []string{"type", "limit_type"})

	// Actieve email batches
	activeEmailBatches := factory.NewGauge(prometheus.GaugeOpts{
		Name: "email_service_active_batches",
		Help: "Het huidige aantal actieve email batches",
	})

	return &PrometheusMetrics{
		emailsSent:           emailsSent,
		emailsFailed:         emailsFailed,
		emailLatency:         emailLatency,
		rateLimitExceeded:    rateLimitExceeded,
		activeEmailBatches:   activeEmailBatches,
		emailTypeCardinality: make(map[string]bool),
	}
}

// NewPrometheusMetrics is een compatibiliteitsfunctie die GetPrometheusMetrics aanroept
// Deprecated: Gebruik GetPrometheusMetrics() in plaats hiervan
func NewPrometheusMetrics() *PrometheusMetrics {
	return GetPrometheusMetrics()
}
