package worker

import (
	"urlchecker/common"
	"urlchecker/worker/config"

	"encoding/json"

	"log"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

type TaskSched struct {
	Tasks map[string]*common.Task
}

func (t *TaskSched) autoSync() {
	for _, location := range config.G_config.Locations {
		key := common.TASK_DIR + location
		common.G_taskMgr.Watch(key, t.handlerWatchEvent)
	}
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

	t.Tasks[task.Id] = &task
	log.Printf("更新任务：%s成功，当前任务个数：%d", task.Id, len(t.Tasks))
}

func (t *TaskSched) deleteTask(kv *mvccpb.KeyValue) {
	taskId := common.ParseTaskId(string(kv.Key))

	_, exists := t.Tasks[taskId]
	if exists {
		delete(t.Tasks, taskId)
		log.Printf("删除任务：%s成功", taskId)
	}
}

func (t *TaskSched) Run() {

}

var (
	G_taskSched *TaskSched
)

func InitTaskSched() (err error) {
	mapTasks := make(map[string]*common.Task)

	var tasks []*common.Task
	tasks, err = common.G_taskMgr.GetTasks(config.G_config.Locations)
	if err != nil {
		return
	}

	for _, task := range tasks {
		mapTasks[task.Id] = task
	}

	G_taskSched = &TaskSched{
		Tasks: mapTasks,
	}

	G_taskSched.autoSync()

	return
}
