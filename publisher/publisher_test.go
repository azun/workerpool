package publisher

import (
	"testing"

	"github.com/azun/worker-pool/common"
)

func TestQueueSomeTasks(t *testing.T) {
	queuer := NewPublisher()
	task := common.Task{
		Work: "doSomething",
	}
	queuer.PublishTask(task)
}

func TestQueueJob(t *testing.T) {
	queuer := NewPublisher()

	defer queuer.Close()
	tasks := []common.Task{
		{
			Work: "doSomething",
		},
	}
	queuer.PublishJob(tasks)

}
