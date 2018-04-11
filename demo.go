package main

import (
	"webserver"
	"webserver/config"
	"webserver/handler"
	"gopkg.in/mgo.v2"
	"webserver/server"
)

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
	args.GetQuery("get query key from url")
	//return type string
	args.GetJsonKey("get Json Key from Post Context")
	//return type *simplejson.Json
	args.GetUserId() //get user Id from session
	//return bool(is user valid) and string(user Id,must be bson.ObjectId.Hex string)
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
	AnyCollection *mgo.Collection
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
	//get new web-server container from your handler and your config
	svr := webserver.NewWebServerFromHandler(new(Config), new(ApiHandler))
	svr.Start()
}
