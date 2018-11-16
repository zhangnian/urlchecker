package worker

import (
	"context"
	"github.com/mongodb/mongo-go-driver/mongo"
	"log"
	"time"
	"urlchecker/common"
	"urlchecker/worker/config"
)

type ResultSaver struct {
	Client     *mongo.Client
	Collection *mongo.Collection
	Q          chan *common.TaskResult
}

func (r *ResultSaver) write(taskResult *common.TaskResult) {
	select {
	case r.Q <- taskResult:
	default:

	}
}

func (r *ResultSaver) writeLoop() {
	var (
		batch []interface{}
	)

	timer := time.NewTicker(time.Second * 10)
	for {
		select {
		case taskResult := <-r.Q:
			if batch == nil {
				log.Println("分配batch空间")
				batch = make([]interface{}, 0)
			}

			batch = append(batch, taskResult)
			if len(batch) >= 100 {
				r.Collection.InsertMany(context.TODO(), batch)
				batch = nil
			}
		case <-timer.C:
			if batch != nil {
				r.Collection.InsertMany(context.TODO(), batch)
				batch = nil
			}
		}
	}
}

var (
	G_resultSaver *ResultSaver
)

func InitResultSaver() (err error) {
	var (
		client *mongo.Client
	)

	log.Printf("开始连接mongodb：%s\n", config.G_config.MongodbUri)
	client, err = mongo.Connect(context.TODO(), config.G_config.MongodbUri)
	if err != nil {
		return
	}

	saver := ResultSaver{
		Q:          make(chan *common.TaskResult, 1000),
		Client:     client,
		Collection: client.Database(config.G_config.MongodbDatabase).Collection(config.G_config.MongodbCollection),
	}

	G_resultSaver = &saver

	go G_resultSaver.writeLoop()
	return
}
