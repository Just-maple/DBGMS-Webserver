package main

import (
	"webserver/handler"
	"webserver/server"
	"webserver/logger"
)

var log = logger.Log

type ApiHandler struct {
	//your custom handler
	*Config
	//must contain Config that implement default config and named Config
	*handler.DefaultApiHandler
	//implement default handler and named DefaultApiHandler
	db *DataBase
	//implement singleton DataBase
	
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
	ret = h.db.AnyCollection.Name
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
	h.db = new(DataBase)
	//init your DataBase
	return h.db
	//should return interface implement server.DB
}
