package main

import (
	"logger"
	"webserver/args"
	"webserver/handler"
	"webserver/user"
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
	log.Debug("Register API from Handler")
	h.ApiGetHandlers.RegisterDefaultAPI("test2", h.ApiTest)
	h.ApiPostHandlers.RegisterDefaultAPI("test3", func(args *args.APIArgs) (ret interface{}, err error) {
		d, err := h.db.AnyCollection.GenerateRawStruct()
		log.Debug(d)
		return d, err
	})
	
	//init default user controller with your user database collection
	var userController = user.InitController(h.db.WXUser)
	h.InjectController(userController)
}

//default api handler
func (h *ApiHandler) ApiTest(args *args.APIArgs) (ret interface{}, err error) {
	queryString := args.Query("get query key from url")
	//return type string
	log.Debug(queryString)
	jsonValue := args.JsonKey("get Json Key from Post Context")
	//return type *jsonx.Json
	log.Debug(jsonValue)
	isValidUser, userId := args.UserId() //get user Id from session
	//return bool(is user valid) and string(user Id,must be bson.ObjectId.Hex string)
	log.Debug(isValidUser, userId)
	
	return
}

func (h *ApiHandler) InitMetaConfig() {
	//implement method InitMetaConfig
	log.Debug("Init Some Other Custom Config")
	//you can handle your extend config here
	//this method will execute after database init
	//and before server start
}

func (h *ApiHandler) NewDataBase() handler.DB {
	//implement method NewDataBase
	log.Debug("New DataBase From Handler")
	//init your DataBase
	h.db = new(DataBase)
	return h.db
	//should return interface implement server.DB
}
