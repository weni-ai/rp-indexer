package indexer

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func StartMetrics() {
	r := chi.NewRouter()

	r.Handle("/metrics", promhttp.Handler())

	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(1)

	// Get metrics port from environment or use default
	metricsPort := os.Getenv("METRICS_PORT")
	if metricsPort == "" {
		metricsPort = "8070"
	}

	// and start serving HTTP
	go func() {
		defer waitGroup.Done()
		addr := fmt.Sprintf(":%s", metricsPort)
		log.Info(fmt.Sprintf("metrics server listening on %s", addr),
			"comp", "server",
			"state", "starting",
		)

		err := http.ListenAndServe(addr, r)
		if err != nil && err != http.ErrServerClosed {
			log.WithFields(log.Fields{
				"error": err,
				"port":  metricsPort,
			}).Error("failed to start metrics server, metrics will not be available")
		}
	}()
}
