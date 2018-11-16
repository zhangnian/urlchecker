package main

import (
	"flag"
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

	err = worker.InitResultSaver()
	common.PanicIfError(err)

	err = worker.InitTaskRunner()
	common.PanicIfError(err)

	err = worker.InitTaskSched()
	common.PanicIfError(err)

	worker.G_taskSched.Sched()
}
