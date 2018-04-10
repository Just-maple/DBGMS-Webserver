package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
)

type DefaultConfig struct {
	MgoDBUrl         string
	ServerAddr       string
	TablePath        string
	SessionSecretKey string
	SessionKey       string
}

func LoadConfigurationFromFile(filename string, config interface{}) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("LoadCfgFromFile ioutil.ReadFile(%s) err(%v)", filename, err)
		return
	}
	err = json.Unmarshal(data, config)
	if err != nil {
		log.Fatalf("LoadCfgFromFile json.Unmarshal err(%v)", err)
		return
	}
	return
}

func LoadConfiguration() (config *DefaultConfig) {
	var configPath = flag.String("f", "./config/config-dev.json", "config path")
	config = new(DefaultConfig)
	LoadConfigurationFromFile(*configPath, &config)
	return
}

func (cfg *DefaultConfig) GetSessionSecretKey() string {
	return cfg.SessionSecretKey
}
func (cfg *DefaultConfig) GetSessionKey() string {
	return cfg.SessionKey
}

func (cfg *DefaultConfig) GetMgoDBUrl() string {
	return cfg.MgoDBUrl
}
func (cfg *DefaultConfig) GetTablePath() string {
	return cfg.TablePath
}
func (cfg *DefaultConfig) GetServerAddr() string {
	return cfg.ServerAddr
}
