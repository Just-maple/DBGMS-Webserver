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
}

type ApiHandler struct {
	*Config
	*handler.DefaultApiHandler
}

func (h *ApiHandler) RegisterAPI() {
	h.ApiGetHandlers.RegisterDefaultAPI("test",h.Test)
}

func (h *ApiHandler)Test(args server.DefaultAPIArgs)(ret interface{},err error){
	return
}


func (h *ApiHandler) InitMetaConfig() {

}

func (h *ApiHandler) NewDataBase() server.DB {
	return new(DataBase)
}

type DataBase struct {
	AnyCollection *mgo.Collection
}

func (db *DataBase) AuthSuperAdminUser(userId string) (bool, bool) {
	return userId == "User is Admin", userId == "User is Super"
}

func (db *DataBase) AuthAdminUser(userId string) bool {
	return userId == "User is Admin"
}

func main() {
	svr := webserver.NewWebServerFromHandler(new(Config), new(ApiHandler))
	svr.Start()
}
