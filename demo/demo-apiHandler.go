package main

import (
	"webserver/handler"
	"webserver/logger"
	"gopkg.in/mgo.v2/bson"
)

var log = logger.Log

type ApiHandler struct {
	//your custom handler
	*Config
	//must contain Config that implement default config
	*handler.DefaultApiHandler
	//must implement default handler
	db *DataBase
	//implement singleton DataBase
	
	MetaData
	//you can define any extend data
}

type MetaData struct {
	String    string
	Int       int
	Bool      bool
	Interface interface{}
	Error     error
}

func (h *ApiHandler) RegisterAPI() {
	//implement method RegisterAPI
	//this method provide Api register
	//and will execute before server start
	
	type test struct {
		Id    bson.ObjectId `bson:"_id"`
		Count int
	}
	//`{
	//	"Test": 1,
	//	"Test2": "string"
	//}`
	h.ApiGetHandlers.RegisterDefaultAPI("test2", h.ApiTest)
	h.ApiPostHandlers.RegisterDefaultAPI("test3", func(args *handler.APIArgs) (ret interface{}, err error) {
		var s = test{}
		err = args.JsonKey("test").Unmarshal(&s)
		var d []test
		err = h.db.AnyCollection.FindAll(nil, &d)
		return d, err
	})
}

func (h *ApiHandler) ApiTest(args *handler.APIArgs) (ret interface{}, err error) {
	queryString := args.Query("get query key from url")
	//return type string
	log.Debug(queryString)
	jsonValue := args.JsonKey("get Json Key from Post Context")
	//return type *simplejson.Json
	log.Debug(jsonValue)
	isValidUser, userId := args.UserId() //get user Id from session
	//return bool(is user valid) and string(user Id,must be bson.ObjectId.Hex string)
	log.Debug(isValidUser, userId)
	
	return
}

func (h *ApiHandler) InitMetaConfig() {
	//implement method InitMetaConfig
	
	//you can handle your extend data here
	//this method will execute after database init
}

func (h *ApiHandler) NewDataBase() handler.DB {
	//implement method NewDataBase
	
	//init your DataBase
	h.db = new(DataBase)
	return h.db
	//should return interface implement server.DB
}
