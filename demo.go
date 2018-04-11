package main

import (
	"webserver"
	"webserver/config"
	"webserver/handler"
	"webserver/server"
	"webserver/logger"
	"webserver/dbx"
)

var log = logger.Log

type Config struct {
	config.DefaultConfig
	//implement default config
}

type ApiHandler struct {
	//your custom handler
	*Config
	//must contain Config that implement default config and named Config
	*handler.DefaultApiHandler
	//implement default handler and named DefaultApiHandler
	
	MetaData struct {
		//you can define any extend data
	}
}

func (h *ApiHandler) RegisterAPI() {
	//implement method RegisterAPI
	
	//this method provide Api register
	//and will execute before server start
	h.ApiGetHandlers.RegisterDefaultAPI("test", h.Test)
}

func (h *ApiHandler) Test(args server.DefaultAPIArgs) (ret interface{}, err error) {
	queryString := args.GetQuery("get query key from url")
	//return type string
	log.Debug(queryString)
	
	jsonValue := args.GetJsonKey("get Json Key from Post Context")
	//return type *simplejson.Json
	log.Debug(jsonValue)
	
	isValidUser, userId := args.GetUserId() //get user Id from session
	//return bool(is user valid) and string(user Id,must be bson.ObjectId.Hex string)
	log.Debug(isValidUser, userId)
	return
}

func (h *ApiHandler) InitMetaConfig() {
	//implement method InitMetaConfig
	h.MetaData = struct{}{
	
	}
	//you can handle your extend data here
	//this method will execute after database init
}

func (h *ApiHandler) NewDataBase() server.DB {
	//implement method NewDataBase
	return new(DataBase)
	//should return interface implement server.DB
}

type DataBase struct {
	//database struct interface implement server.DB
	AnyCollection *dbx.Collection
	//any public *mgo.Collection will init when database init
	//collection will init from Collection name in lower case like "anycollection"
}

func (db *DataBase) AuthSuperAdminUser(userId string) (bool, bool) {
	//database struct implement auth super admin user
	//define your logic here
	return userId == "User is Admin", userId == "User is Super"
}

func (db *DataBase) AuthAdminUser(userId string) bool {
	//database struct implement auth  admin user
	//define your logic here
	return userId == "User is Admin"
}

func main() {
	demoConfg := &Config{
		config.DefaultConfig{
			MgoDBUrl:         "mongodb://wx2.asoapp.com:27777/wx",
			ServerAddr:       "0.0.0.0:8888",
			TablePath:        "./",
			SessionSecretKey: "secret",
			SessionKey:       "session",
		},
	}
	//get new web-server container from your handler and your config
	svr := webserver.NewWebServerFromHandler(demoConfg, new(ApiHandler))
	svr.Start()
}
