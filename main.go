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

func (h *ApiHandler) RegisterGETAPI() {

}
func (h *ApiHandler) RegisterPOSTAPI() {

}
func (h *ApiHandler) RegisterSpecificAPI() {

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
	webserver.NewWebServerFromHandler(new(Config), new(ApiHandler))
}
