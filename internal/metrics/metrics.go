package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// QuotesFetchedTotal counts the total number of quotes fetched
	QuotesFetchedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "quotes_fetched_total",
		Help: "Total number of quotes successfully fetched from OpenRouter",
	})

	// QuotesByTag counts quotes by tag
	QuotesByTag = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "quotes_by_tag",
		Help: "Number of quotes fetched by tag",
	}, []string{"tag"})

	// QuoteFetchErrorsTotal counts the total number of errors fetching quotes
	QuoteFetchErrorsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "quote_fetch_errors_total",
		Help: "Total number of errors while fetching quotes",
	})

	// QuoteFetchLatency measures the latency of quote fetch operations
	QuoteFetchLatency = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "quote_fetch_latency_seconds",
		Help:    "Latency of quote fetch operations in seconds",
		Buckets: prometheus.DefBuckets,
	})

	// HTTPRequestsTotal counts HTTP requests by method, route, and status
	HTTPRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"method", "route", "status"})

	// OpenRouterUp indicates if OpenRouter API is up (1) or down (0)
	OpenRouterUp = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "openrouter_up",
		Help: "Indicates if the last OpenRouter API call succeeded (1) or failed (0)",
	})
)

// Init initializes metrics (called at startup)
func Init() {
	// Set initial OpenRouter status to unknown (0)
	OpenRouterUp.Set(0)
}

// RecordQuoteFetched increments the quotes fetched counter
func RecordQuoteFetched(tag string) {
	QuotesFetchedTotal.Inc()
	QuotesByTag.WithLabelValues(tag).Inc()
}

// RecordQuoteError increments the error counter
func RecordQuoteError() {
	QuoteFetchErrorsTotal.Inc()
}

// RecordLatency records the latency of a quote fetch
func RecordLatency(seconds float64) {
	QuoteFetchLatency.Observe(seconds)
}

// RecordHTTPRequest records an HTTP request
func RecordHTTPRequest(method, route, status string) {
	HTTPRequestsTotal.WithLabelValues(method, route, status).Inc()
}

// SetOpenRouterStatus sets the OpenRouter status
func SetOpenRouterStatus(up bool) {
	if up {
		OpenRouterUp.Set(1)
	} else {
		OpenRouterUp.Set(0)
	}
}
