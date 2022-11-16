package worker

import (
	"context"
	"fmt"

	"github.com/azun/worker-pool/common"
	redis "github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type Worker struct {
	ctx context.Context
	rdb *redis.Client
}

func NewWorker() Worker {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return Worker{
		ctx: context.Background(),
		rdb: rdb,
	}
}

func (w Worker) InitRedisSearch() {

}

func (w Worker) Work() {
	result, err := w.rdb.BRPop(w.ctx, 0, "ui_tasks", "import_tasks").Result()
	if err != nil {
		logrus.Error(err)
		panic(err)
	}
	q, taskId := result[0], result[1]
	var task common.Task

	logrus.Debug("Got task from " + q + ". Fetching " + taskId)
	if err := w.rdb.HGetAll(w.ctx, taskId).Scan(&task); err != nil {
		logrus.Error(err)
		panic(err)
	}

	w.fakeWork(task)
	logrus.Debug("Finished task " + task.Id + ". Cleanup")
	if task.Job != "" {
		jobTasksPending, err := w.rdb.HIncrBy(w.ctx, task.Job, "numTasks", -1).Result()
		if err != nil {
			logrus.Error(err)
			panic(err)
		}
		logrus.Debug("Job " + task.Job + " has " + fmt.Sprint(jobTasksPending) + " tasks pending")
		if jobTasksPending <= 0 {
			w.rdb.Del(w.ctx, task.Job)
		}

	}
	w.rdb.Del(w.ctx, taskId)
}

func (w Worker) fakeWork(task common.Task) {
	logrus.Debug("Processing task " + task.Work)
}
