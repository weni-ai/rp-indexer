package indexer

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func StartMetrics() {
	r := chi.NewRouter()

	r.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8070", r)
}
