package main

import (
	"flag"
	"time"
	"urlchecker/common"
	"urlchecker/master/config"
	"urlchecker/master/http"
)

var (
	cfg = flag.String("config", "./master.json", "config file path")
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

	http.Init()

	time.Sleep(time.Hour * 1)
}
