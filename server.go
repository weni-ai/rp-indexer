package indexer

import (
	"fmt"
	"net/http"
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

	// and start serving HTTP
	go func() {
		defer waitGroup.Done()
		err := http.ListenAndServe(":8070", r)
		if err != nil && err != http.ErrServerClosed {
			log.Error("failed to start server", "error", err, "comp", "server", "state", "stopping")
		}
	}()

	log.Info(fmt.Sprintf("server listening on 8070"),
		"comp", "server",
		"state", "started",
	)
}
