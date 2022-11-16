package main

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/azun/worker-pool/cmd/params"
	"github.com/azun/worker-pool/common"
	"github.com/azun/worker-pool/publisher"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var (
	tasksPublished = promauto.NewCounter(prometheus.CounterOpts{
		Name: "publisher_tasks_total",
		Help: "The total number of published tasks",
	})
	jobsPublished = promauto.NewCounter(prometheus.CounterOpts{
		Name: "publisher_jobs_total",
		Help: "The total number of published jobs",
	})
)

func main() {

	var conf params.Params
	conf.ParseArgs()

	logrus.Info("Starting publisher")
	go publishJobs(conf)

	logrus.Info("Starting metrics server")
	http.Handle("/metrics", promhttp.Handler())
	logrus.Error(http.ListenAndServe(":2020", nil))
}

func publishJobs(conf params.Params) {
	p := publisher.NewPublisher()

	for i := 0; i != conf.Testing.NumJobs; i++ {
		numTasks := 1 + rand.Intn(conf.Testing.JobMaxTasks)
		var tasks []common.Task
		for j := 0; j < numTasks; j++ {
			tasks = append(tasks, common.Task{
				Work: "doSomething - " + fmt.Sprint(j),
			})
		}
		p.PublishJob(tasks)
		jobsPublished.Inc()
		tasksPublished.Add(float64(numTasks))
	}
}
