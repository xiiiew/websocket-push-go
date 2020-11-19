package common

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
	"sync"
)

type config struct {
	HttpWsPort    string
	RpcPort       string
	PingDuration  int
	CloseDuration int
	InChanLength  int
	OutChanLength int
}

var Config *config

func initConfig() {
	HttpWsPort := configFile.Section("application").Key("http-ws-port").String()
	RpcPort := configFile.Section("application").Key("rpc-port").String()
	PingDuration, _ := configFile.Section("ws-server").Key("ping_duration").Int()
	CloseDuration, _ := configFile.Section("ws-server").Key("close_duration").Int()
	InChanLength, _ := configFile.Section("ws-server").Key("in_chan_length").Int()
	OutChanLength, _ := configFile.Section("ws-server").Key("out_chan_length").Int()

	Config = &config{
		HttpWsPort:    HttpWsPort,
		RpcPort:       RpcPort,
		PingDuration:  PingDuration,
		CloseDuration: CloseDuration,
		InChanLength:  InChanLength,
		OutChanLength: OutChanLength,
	}
}

var once sync.Once
var configFile *ini.File

// 单例加载配置文件
func init() {
	once.Do(loadConfig)
}

func loadConfig() {
	pathConfig := "config/config.ini"
	cfg, err := ini.Load(pathConfig)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	configFile = cfg
	initConfig()
}
