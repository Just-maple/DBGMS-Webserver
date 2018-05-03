package main

import "webserver/config"

type Config struct {
	config.DefaultConfig
	//implement default config
}

func NewConfig() *Config {
	return &Config{
		config.DefaultConfig{
			MgoDBUrl: "mongodb://localhost:27017/test",
			//ServerAddr:       "0.0.0.0:8388",
			//TablePath:        "./table/",
			//SessionSecretKey: "secret",
			//SessionKey:       "session",
		},
	}
}
