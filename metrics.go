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
	esIndexingTimeSummary.WithLabelValues(operation).Observe(duration)
}
