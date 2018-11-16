package common

import (
	"log"
	"time"

	"fmt"

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

func InitTaskMgr(ectdHost string, etcdPort int) (err error) {
	cfg := clientv3.Config{
		Endpoints:   []string{fmt.Sprintf("%s:%d", ectdHost, etcdPort)},
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

func (t *TaskMgr) SaveTask(task *Task) (err error) {
	var v []byte

	for _, location := range task.Locations {
		key := TASK_DIR + location + "/" + task.Id
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

func (t *TaskMgr) GetTasks(locations []string) (tasks []*Task, err error) {
	tasks = make([]*Task, 0)

	var getResp *clientv3.GetResponse
	for _, location := range locations {
		key := TASK_DIR + location
		getResp, err = t.Kv.Get(context.TODO(), key, clientv3.WithPrefix())
		if err != nil {
			return
		}

		for _, keyvalue := range getResp.Kvs {
			var task Task
			err = json.Unmarshal(keyvalue.Value, &task)
			if err != nil {
				err = nil
				continue
			}

			tasks = append(tasks, &task)
		}
	}

	return
}

func (t *TaskMgr) Watch(key string, handler func(clientv3.WatchResponse)) (err error) {
	watchChan := t.Client.Watch(context.TODO(), key, clientv3.WithPrefix(), clientv3.WithPrevKV())
	log.Printf("watch key: %s", key)

	go func() {
		for {
			select {
			case watchResp := <-watchChan:
				handler(watchResp)
			}
		}
	}()

	return
}
