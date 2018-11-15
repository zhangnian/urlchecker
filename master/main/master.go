package main

import (
	"flag"
	"time"
	"urlchecker/master"
	"urlchecker/master/config"
	"urlchecker/master/http"
)

var (
	cfg = flag.String("config", "./master.json", "config file path")
)

func initFlag() {
	flag.Parse()
}

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	initFlag()

	err := config.Init(*cfg)
	PanicIfError(err)

	err = master.InitTaskMgr()
	PanicIfError(err)

	http.Init()

	time.Sleep(time.Hour * 1)
}
