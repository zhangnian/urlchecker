package main

import (
	"flag"
	"time"
	"urlchecker/common"
	"urlchecker/worker"
	"urlchecker/worker/config"
)

var (
	cfg = flag.String("config", "./worker.json", "config file path")
)

func initFlag() {
	flag.Parse()
}

func main() {
	initFlag()

	err := config.Init(*cfg)
	common.PanicIfError(err)

	err = common.InitTaskMgr(config.G_config.EtcdHost, config.G_config.EtcdPort)
	common.PanicIfError(err)

	err = worker.InitTaskSched()
	common.PanicIfError(err)

	time.Sleep(time.Hour * 1)
}
