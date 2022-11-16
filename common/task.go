package common

import "encoding/json"

type Task struct {
	Id   string
	Work string
	Job  string
}

func (t *Task) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Task) UnmarshalBinary(data string) error {
	if err := json.Unmarshal([]byte(data), &t); err != nil {
		return err
	}
	return nil
}
