package config

import (
	"encoding/json"
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

func LoadConfiguration(configPath string) (config *DefaultConfig) {
	//var configPath = flag.String("f", "./config/config-dev.json", "config path")
	config = new(DefaultConfig)
	LoadConfigurationFromFile(configPath, &config)
	return
}

func (cfg *DefaultConfig) GetSessionSecretKey() string {
	if cfg.SessionSecretKey == "" {
		cfg.SessionSecretKey = "secret-key"
	}
	return cfg.SessionSecretKey
}
func (cfg *DefaultConfig) GetSessionKey() string {
	if cfg.SessionKey == "" {
		cfg.SessionKey = "session"
	}
	return cfg.SessionKey
}

func (cfg *DefaultConfig) GetMgoDBUrl() string {
	if cfg.MgoDBUrl == "" {
		cfg.MgoDBUrl = "0.0.0.0:27017"
	}
	return cfg.MgoDBUrl
}
func (cfg *DefaultConfig) GetTablePath() string {
	if cfg.TablePath == "" {
		cfg.TablePath = "./"
	}
	return cfg.TablePath
}
func (cfg *DefaultConfig) GetServerAddr() string {
	if cfg.ServerAddr == "" {
		cfg.ServerAddr = "0.0.0.0:8388"
	}
	return cfg.ServerAddr
}
