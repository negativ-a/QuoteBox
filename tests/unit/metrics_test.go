package unit

import (
	"testing"

	"github.com/Adeel56/quotebox/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestRecordQuoteFetched(t *testing.T) {
	// Reset metrics before test
	metrics.QuotesFetchedTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "quotes_fetched_total_test",
	})
	metrics.QuotesByTag = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "quotes_by_tag_test",
	}, []string{"tag"})

	tag := "joy"
	metrics.RecordQuoteFetched(tag)

	// Verify counter was incremented
	count := testutil.ToFloat64(metrics.QuotesFetchedTotal)
	assert.Equal(t, float64(1), count)
}

func TestRecordQuoteError(t *testing.T) {
	// Reset metrics before test
	metrics.QuoteFetchErrorsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "quote_fetch_errors_total_test",
	})

	metrics.RecordQuoteError()

	count := testutil.ToFloat64(metrics.QuoteFetchErrorsTotal)
	assert.Equal(t, float64(1), count)
}

func TestSetOpenRouterStatus(t *testing.T) {
	// Reset metric before test
	metrics.OpenRouterUp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "openrouter_up_test",
	})

	// Test setting to up
	metrics.SetOpenRouterStatus(true)
	assert.Equal(t, float64(1), testutil.ToFloat64(metrics.OpenRouterUp))

	// Test setting to down
	metrics.SetOpenRouterStatus(false)
	assert.Equal(t, float64(0), testutil.ToFloat64(metrics.OpenRouterUp))
}

func TestRecordLatency(t *testing.T) {
	// Reset metric before test
	metrics.QuoteFetchLatency = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "quote_fetch_latency_seconds_test",
		Buckets: prometheus.DefBuckets,
	})

	// Register the histogram for testing
	reg := prometheus.NewRegistry()
	reg.MustRegister(metrics.QuoteFetchLatency)

	metrics.RecordLatency(0.5)
	metrics.RecordLatency(1.0)
	metrics.RecordLatency(2.0)

	// Verify histogram has observations by checking the metric family
	metricFamilies, err := reg.Gather()
	assert.NoError(t, err)
	assert.Len(t, metricFamilies, 1)
	
	// Check the histogram count
	histogram := metricFamilies[0].GetMetric()[0].GetHistogram()
	assert.Equal(t, uint64(3), histogram.GetSampleCount())
	assert.Equal(t, float64(3.5), histogram.GetSampleSum())
}
