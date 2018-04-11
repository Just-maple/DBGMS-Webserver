package main

import "webserver/config"

type Config struct {
	config.DefaultConfig
	//implement default config
}

func NewConfig() *Config {
	return &Config{
		config.DefaultConfig{
			MgoDBUrl:         "mongodb://wx2.asoapp.com:27777/wx",
			ServerAddr:       "0.0.0.0:8388",
			TablePath:        "./",
			SessionSecretKey: "secret",
			SessionKey:       "session",
		},
	}
}
