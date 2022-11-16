package publisher

import (
	"context"
	"fmt"

	"github.com/azun/worker-pool/common"
	redis "github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jeremywohl/flatten"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

type Publisher struct {
	ctx context.Context
	rdb *redis.Client
}

func NewPublisher() Publisher {
	return Publisher{
		ctx: context.Background(),
		rdb: redis.NewClient(&redis.Options{
			Addr:     "redis:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
	}
}

func (q Publisher) PublishTask(task common.Task) {
	pipe := q.rdb.Pipeline()

	task.Id = "task:" + uuid.NewString()

	var payload map[string]interface{}

	mapstructure.Decode(task, &payload)

	logrus.Debug("Saving task " + task.Id)
	if err := pipe.HSet(q.ctx, task.Id, payload).Err(); err != nil {
		logrus.Error(err)
		panic(err)
	}
	logrus.Debug("Queueing task " + task.Id)
	if err := pipe.LPush(q.ctx, "import_tasks", task.Id).Err(); err != nil {
		logrus.Error(err)
		panic(err)
	}

	pipe.Exec(q.ctx)

}
func (q Publisher) PublishJob(tasks []common.Task) {
	_, err := q.rdb.Pipelined(q.ctx, func(p redis.Pipeliner) error {
		job := common.NewJob(tasks)

		logrus.Debugf("Creating %s with %d tasks", job.Id, len(tasks))

		payload, err := flatten.Flatten(job.ToMap(), "", flatten.DotStyle)

		err = p.HSet(q.ctx, job.Id, payload).Err()
		if err != nil {
			logrus.Error(err)
			panic(err)
		}

		for i, _ := range tasks {
			logrus.Debugf("Queueing task %d", i)
			taskId := fmt.Sprintf("%s.tasks.%d", job.Id, i)
			if err := p.LPush(q.ctx, "import_tasks", taskId).Err(); err != nil {
				logrus.Error(err)
				panic(err)
			}
		}

		logrus.Debug("Sending job " + job.Id)

		return nil

	})
	if err != nil {
		logrus.Error(err)
		panic(err)
	}
}

func (q Publisher) Close() {
	logrus.Debug("Closing redis connection")
	q.rdb.Pipeline().Close()
	q.rdb.Close()
}
