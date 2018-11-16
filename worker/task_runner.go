package worker

import (
	"log"
	"time"
	"urlchecker/common"
	"urlchecker/worker/config"
)

type TaskRunner struct {
	Q               chan *common.TaskPlan
	ConcurrencyChan chan struct{}
}

func (t *TaskRunner) RunTask(taskPlan *common.TaskPlan) (err error) {
	taskPlan.IsSched = true

	select {
	case t.Q <- taskPlan:
	default:
	}

	return
}

func (t *TaskRunner) runLoop() {
	for {
		select {
		case taskPlan := <-t.Q:
			go t.executeTask(taskPlan)
		}
	}
}

func (t *TaskRunner) executeTask(taskPlan *common.TaskPlan) {
	log.Printf("开始执行任务：%s\n", taskPlan.Task.Id)
	var (
		statusCode    int
		contentLength int64
		err           error
	)

	runAt := time.Now()
	if taskPlan.Task.Method == "GET" {
		statusCode, contentLength, err = HttpGet(taskPlan.Task.Uri)
	}
	finishAt := time.Now()

	taskPlan.IsSched = false

	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	taskResult := common.TaskResult{
		TaskId:        taskPlan.Task.Id,
		Uri:           taskPlan.Task.Uri,
		Method:        taskPlan.Task.Method,
		Err:           errMsg,
		StatusCode:    statusCode,
		ContentLength: contentLength,
		RunAt:         runAt,
		FinishAt:      finishAt,
		Cost:          finishAt.Sub(runAt).Nanoseconds() / 1000 / 1000,
	}
	G_taskSched.PushResult(&taskResult)
}

var (
	G_taskRunner *TaskRunner
)

func InitTaskRunner() (err error) {
	taskRunner := TaskRunner{
		Q:               make(chan *common.TaskPlan, 1000),
		ConcurrencyChan: make(chan struct{}, config.G_config.MaxConcurrency),
	}

	G_taskRunner = &taskRunner

	go G_taskRunner.runLoop()
	return
}
