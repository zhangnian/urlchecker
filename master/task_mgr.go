package master

import (
	"time"

	"fmt"
	"urlchecker/master/config"

	"urlchecker/common"

	"context"

	"encoding/json"

	"go.etcd.io/etcd/clientv3"
)

type TaskMgr struct {
	Client *clientv3.Client
	Kv     clientv3.KV
}

var (
	G_taskMgr *TaskMgr
)

func InitTaskMgr() (err error) {
	cfg := clientv3.Config{
		Endpoints:   []string{fmt.Sprintf("%s:%d", config.G_config.EtcdHost, config.G_config.EtcdPort)},
		DialTimeout: 5 * time.Second,
	}
	client, err := clientv3.New(cfg)
	if err != nil {
		return
	}

	taskMgr := TaskMgr{
		Client: client,
		Kv:     client.KV,
	}

	G_taskMgr = &taskMgr
	return
}

func (t *TaskMgr) SaveTask(task *common.Task) (err error) {
	var v []byte

	for _, location := range task.Locations {
		key := common.TASK_DIR + location + "/" + task.Id
		v, err = json.Marshal(task)
		if err != nil {
			return
		}
		t.Kv.Put(context.TODO(), key, string(v))
	}

	return
}

func (t *TaskMgr) DeleteTask(taskIds []string) (err error) {
	return
}
