package main

import (
	"net/http"
	"runtime"
	"sync"

	"github.com/azun/worker-pool/worker"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "worker_messages_total",
		Help: "The total number of processed events",
	})
)

func main() {
	// logrus.SetLevel(logrus.DebugLevel)

	logrus.Info("Starting worker app")
	w := worker.NewWorker()
	w.InitRedisSearch()
	wg := sync.WaitGroup{}

	numCpu := 2 * runtime.NumCPU()
	logrus.Infof("Spawning %d threads", numCpu)
	wg.Add(numCpu)
	for i := 0; i < numCpu; i++ {
		go runWorkerThread(w, &wg)
	}

	logrus.Info("Starting metrics server")
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2021", nil)
}

func runWorkerThread(w worker.Worker, wg *sync.WaitGroup) {
	defer wg.Done()
	logrus.Debug("Starting worker. Listening on queues")

	for {
		w.Work()
		opsProcessed.Inc()
	}
}
