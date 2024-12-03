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

var contactsInQueue = promauto.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "indexer_contacts_in_queue",
		Help: "Number of contacts currently in the queue waiting to be indexed.",
	},
	[]string{"queue"},
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

func UpdateContactsInQueue(queueName string, count int) {
	contactsInQueue.WithLabelValues(queueName).Set(float64(count))
}

func ObserveDBResponseTime(operation string, duration float64) {
	dbResponseTimeSummary.WithLabelValues(operation).Observe(duration)
}

func ObserveESIndexingTime(index string, duration float64) {
	esIndexingTimeSummary.WithLabelValues(index).Observe(duration)
}
