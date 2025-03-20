package indexer

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var summaryObjectives = map[float64]float64{
	0.5:  0.05,  // 50th percentile with a max. absolute error of 0.05.
	0.90: 0.01,  // 90th percentile with a max. absolute error of 0.01.
	0.95: 0.005, // 95th percentile with a max. absolute error of 0.005.
	0.99: 0.001, // 99th percentile with a max. absolute error of 0.001.
}

var contactsProcessing = promauto.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "indexer_contacts_processing",
		Help: "Number of contacts created / deleted / conflited per batch (500)",
	},
	[]string{"process"},
)

var dbResponseTimeSummary = promauto.NewSummaryVec(
	prometheus.SummaryOpts{
		Name:       "db_response_time_summary_seconds",
		Help:       "Database response time in seconds.",
		Objectives: summaryObjectives,
	},
	[]string{"operation"},
)

var esIndexingTimeSummary = promauto.NewSummaryVec(
	prometheus.SummaryOpts{
		Name:       "es_indexing_time_summary_seconds",
		Help:       "Indexing time in Elasticsearch in seconds.",
		Objectives: summaryObjectives,
	},
	[]string{"operation"},
)

var elapsedTimeSinceIndexing = promauto.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "elapsed_time_since_indexing",
		Help: "Total elapsed time of indexing (500k)",
	},
	[]string{"operation"},
)

// New metrics below
var totalContactsProcessed = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "indexer_total_contacts_processed",
		Help: "Total number of contacts processed since start",
	},
	[]string{"process"},
)

var indexingErrors = promauto.NewCounter(
	prometheus.CounterOpts{
		Name: "indexer_errors_total",
		Help: "Total number of errors encountered during indexing",
	},
)

var indexingLatency = promauto.NewHistogram(
	prometheus.HistogramOpts{
		Name:    "indexer_contact_latency_seconds",
		Help:    "Time between contact modification and successful indexing",
		Buckets: prometheus.ExponentialBuckets(0.1, 2, 10), // From 0.1s to ~102s
	},
)

func UpdateContactsPerBatch(process string, count int) {
	contactsProcessing.WithLabelValues(process).Set(float64(count))
}

func ObserveDBResponseTime(operation string, duration float64) {
	dbResponseTimeSummary.WithLabelValues(operation).Observe(duration)
}

func ObserveESIndexingTime(operation string, duration float64) {
	esIndexingTimeSummary.WithLabelValues(operation).Observe(duration)
}

func ObserveElapsedIndexingTime(operation string, duration float64) {
	elapsedTimeSinceIndexing.WithLabelValues(operation).Set(duration)
}

// New functions for the new metrics
func IncrementTotalContacts(process string, count int) {
	totalContactsProcessed.WithLabelValues(process).Add(float64(count))
}

func IncrementIndexingErrors() {
	indexingErrors.Inc()
}

func ObserveIndexingLatency(duration float64) {
	indexingLatency.Observe(duration)
}
