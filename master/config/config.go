package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	EtcdHost string `json:"EtcdHost"`
	EtcdPort int    `json:"EtcdPort"`
	ApiHost  string `json:"ApiHost"`
	ApiPort  int    `json:"ApiPort"`
}

var (
	G_config *Config
)

func Init(cfg string) (err error) {
	var b []byte
	if b, err = ioutil.ReadFile(cfg); err != nil {
		return
	}

	var config Config
	if err = json.Unmarshal(b, &config); err != nil {
		return
	}

	G_config = &config

	return
}
