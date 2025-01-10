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
		Help: "Number of contacts created / deleted / conflited per batch",
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
	[]string{"index"},
)

func UpdateContactsPerBatch(processName string, count int) {
	contactsProcessing.WithLabelValues(processName).Set(float64(count))
}

func ObserveDBResponseTime(operation string, duration float64) {
	dbResponseTimeSummary.WithLabelValues(operation).Observe(duration)
}

func ObserveESIndexingTime(index string, duration float64) {
	esIndexingTimeSummary.WithLabelValues(index).Observe(duration)
}
