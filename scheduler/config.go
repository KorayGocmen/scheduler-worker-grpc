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

type httpServerConfig struct {
	Addr string `toml:"addr"`
}

type tomlConfig struct {
	GRPCServer grpcServerConfig `toml:"grpc_server"`
	HTTPServer httpServerConfig `toml:"http_server"`
}

var (
	config tomlConfig
)

func loadConfig() {
	var lConfig tomlConfig
	if _, err := toml.DecodeFile("scheduler/config.toml", &lConfig); err != nil {
		log.Fatalln("Error decoding config file", err)
	}

	config = lConfig
}
