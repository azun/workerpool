package common

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Job struct {
	Id        string
	NumTasks  int
	CreatedAt int64
	Tasks     []Task
}

func NewJob(tasks []Task) Job {
	return Job{
		Id:        "job:" + uuid.NewString(),
		NumTasks:  len(tasks),
		Tasks:     tasks,
		CreatedAt: time.Now().Unix(),
	}
}

func (j *Job) ToMap() map[string]interface{} {
	bytes, err := j.MarshalBinary()
	if err != nil {
		panic(err)
	}

	var payload map[string]interface{}
	err = json.Unmarshal(bytes, &payload)
	if err != nil {
		panic(err)
	}
	return payload
}

func (j *Job) MarshalBinary() ([]byte, error) {
	return json.Marshal(j)
}

func (j *Job) UnmarshalBinary(data string) error {
	if err := json.Unmarshal([]byte(data), &j); err != nil {
		return err
	}
	return nil
}
