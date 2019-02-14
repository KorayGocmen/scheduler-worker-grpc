package main

import (
	"log"

	"github.com/BurntSushi/toml"
)

type grpcServerConfig struct {
	Addr    string `toml:"addr"`
	UseTLS  bool   `toml:"use_tls"`
	CrtFile string `toml:"crt_file"`
	KeyFile string `toml:"key_file"`
}

type schedulerConfig struct {
	Addr string `toml:"addr"`
}

type tomlConfig struct {
	GRPCServer grpcServerConfig `toml:"grpc_server"`
	Scheduler  schedulerConfig  `toml:"scheduler"`
}

var (
	config tomlConfig
)

func loadConfig() {
	var localConfig tomlConfig
	if _, err := toml.DecodeFile("worker/config.toml", &localConfig); err != nil {
		log.Fatalln("Error decoding config file", err)
	}

	config = localConfig
}
