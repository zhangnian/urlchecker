package common

import (
	"fmt"
	"strings"
	"time"

	"github.com/gorhill/cronexpr"
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

type TaskPlan struct {
	Task     *Task
	CronExpr *cronexpr.Expression
	RunAt    time.Time
	IsSched  bool
}

type TaskResult struct {
	TaskId        string    `bson:"taskId"`
	Uri           string    `bson:"uri"`
	Method        string    `bson:"method"`
	Err           string    `bson:"err"`
	StatusCode    int       `bson:"statusCode"`
	ContentLength int64     `bson:"contentLength"`
	RunAt         time.Time `bson:"runAt"`
	FinishAt      time.Time `bson:"finishAt"`
	Cost          int64     `bson:"cost"`
}
