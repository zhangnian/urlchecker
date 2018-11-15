package common

import (
	"fmt"

	"strings"

	"github.com/satori/go.uuid"
)

type Task struct {
	Id        string   `json:"id"`
	Name      string   `json:"name"`
	Cron      string   `json:"cron"`
	Uri       string   `json:"uri"`
	Method    string   `json:"method"`
	Locations []string `json:"locations"`
}

func (t *Task) NewId() (err error) {
	var id uuid.UUID
	if id, err = uuid.NewV4(); err != nil {
		return
	}

	taskId := fmt.Sprintf("%s", id)
	t.Id = strings.Replace(taskId, "-", "", -1)
	return
}
