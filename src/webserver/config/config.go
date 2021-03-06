package config

import (
	"encoding/json"
	"io/ioutil"
	"logger"
)

var log = logger.Log

type DefaultConfig struct {
	MgoDBUrl           string
	ServerAddr         string
	TablePath          string
	SessionSecretKey   string
	SessionKey         string
	SessionExpiredTime int
}

const (
	defaultSecretKey          = "secret-key"
	defaultSessionKey         = "session"
	defaultMgoUrl             = "mongodb://0.0.0.0:27017"
	defaultTablePath          = "./table/"
	defaultServerAddr         = "0.0.0.0:8388"
	defaultSessionExpiredTime = 2 * 24 * 60 * 60 //two days
)

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
		cfg.SessionSecretKey = defaultSecretKey
		log.Noticef("Undefined SessionSecretKey,use default [ %v ]", defaultSecretKey)
	}
	return cfg.SessionSecretKey
}

func (cfg *DefaultConfig) GetSessionKey() string {
	if cfg.SessionKey == "" {
		cfg.SessionKey = defaultSessionKey
		log.Noticef("Undefined SessionKey,use default [ %v ]", defaultSessionKey)
	}
	return cfg.SessionKey
}

func (cfg *DefaultConfig) GetMgoDBUrl() string {
	if cfg.MgoDBUrl == "" {
		cfg.MgoDBUrl = defaultMgoUrl
		log.Noticef("Undefined MgoDBUrl,use default [ %v ]", defaultMgoUrl)
	}
	return cfg.MgoDBUrl
}

func (cfg *DefaultConfig) GetTablePath() string {
	if cfg.TablePath == "" {
		cfg.TablePath = defaultTablePath
		log.Noticef("Undefined TablePath,use default [ %v ]", defaultTablePath)
	}
	return cfg.TablePath
}

func (cfg *DefaultConfig) GetSessionExpiredTime() int {
	if cfg.SessionExpiredTime == 0 {
		cfg.SessionExpiredTime = defaultSessionExpiredTime
		log.Noticef("Undefined SessionExpiredTime,use default [ %v ]",
			defaultSessionExpiredTime)
	}
	return cfg.SessionExpiredTime
}
func (cfg *DefaultConfig) GetServerAddr() string {
	if cfg.ServerAddr == "" {
		cfg.ServerAddr = defaultServerAddr
		log.Noticef("Undefined ServerAddr,use default [ %v ]", defaultServerAddr)
	}
	return cfg.ServerAddr
}
