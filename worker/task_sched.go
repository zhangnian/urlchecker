package worker

import (
	"github.com/gorhill/cronexpr"
	"time"
	"urlchecker/common"
	"urlchecker/worker/config"

	"encoding/json"

	"log"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

type TaskSched struct {
	TaskPlans  map[string]*common.TaskPlan
	TaskResult chan *common.TaskResult
}

func (t *TaskSched) autoSync() {
	for _, location := range config.G_config.Locations {
		key := common.TASK_DIR + location
		common.G_taskMgr.Watch(key, t.handlerWatchEvent)
	}
}

func (t *TaskSched) writeResult() {
	go func() {
		for taskResult := range t.TaskResult {
			G_resultSaver.write(taskResult)
		}
	}()
}

func (t *TaskSched) handlerWatchEvent(watchResp clientv3.WatchResponse) {
	for _, event := range watchResp.Events {
		switch event.Type {
		case mvccpb.PUT:
			t.updateTask(event.Kv)
		case mvccpb.DELETE:
			t.deleteTask(event.Kv)
		}
	}
}

func (t *TaskSched) updateTask(kv *mvccpb.KeyValue) {
	var task common.Task
	err := json.Unmarshal(kv.Value, &task)
	if err != nil {
		return
	}

	cronExpr, err := cronexpr.Parse(task.Cron)
	if err != nil {
		return
	}

	t.TaskPlans[task.Id] = &common.TaskPlan{
		Task:     &task,
		CronExpr: cronExpr,
		RunAt:    cronExpr.Next(time.Now()),
	}
	log.Printf("更新任务：%s成功，当前任务个数：%d", task.Id, len(t.TaskPlans))
}

func (t *TaskSched) deleteTask(kv *mvccpb.KeyValue) {
	taskId := common.ParseTaskId(string(kv.Key))

	_, exists := t.TaskPlans[taskId]
	if exists {
		delete(t.TaskPlans, taskId)
		log.Printf("删除任务：%s成功，当前任务个数：%d", taskId, len(t.TaskPlans))
	}
}

func (t *TaskSched) PushResult(taskResult *common.TaskResult) {
	select {
	case t.TaskResult <- taskResult:
	default:
	}
}

func (t *TaskSched) Sched() {
	for {
		var nearTime *time.Time
		now := time.Now()
		for _, taskPlan := range t.TaskPlans {
			if taskPlan.RunAt.Before(now) || taskPlan.RunAt.Equal(now) {
				if taskPlan.IsSched{
					continue
				}

				G_taskRunner.RunTask(taskPlan)
				taskPlan.RunAt = taskPlan.CronExpr.Next(now)
			}

			if nearTime == nil || taskPlan.RunAt.Before(*nearTime){
				nearTime = &taskPlan.RunAt
			}
		}

		//time.Sleep(time.Millisecond * 100)
		sleepMS := (*nearTime).Sub(time.Now()).Nanoseconds() / 1000 / 1000
		log.Printf("sleep: %dms\n", sleepMS)
		time.Sleep(time.Duration(sleepMS) * time.Millisecond)
	}

}

var (
	G_taskSched *TaskSched
)

func InitTaskSched() (err error) {
	mapTasks := make(map[string]*common.TaskPlan)

	var tasks []*common.Task
	tasks, err = common.G_taskMgr.GetTasks(config.G_config.Locations)
	if err != nil {
		return
	}

	now := time.Now()
	for _, task := range tasks {
		cronExpr, err := cronexpr.Parse(task.Cron)
		if err != nil {
			log.Println("解析任务的cron表达式失败")
			continue
		}

		mapTasks[task.Id] = &common.TaskPlan{
			Task:     task,
			CronExpr: cronExpr,
			RunAt:    cronExpr.Next(now),
		}
	}

	G_taskSched = &TaskSched{
		TaskPlans:  mapTasks,
		TaskResult: make(chan *common.TaskResult, 1000),
	}

	G_taskSched.autoSync()
	G_taskSched.writeResult()

	return
}
