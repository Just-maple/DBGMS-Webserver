package main

import (
	"webserver/handler"
	"webserver/logger"
	"webserver/user"
	"webserver/permission"
	"reflect"
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
	
	h.ApiGetHandlers.RegisterDefaultAPI("test2", h.ApiTest)
	h.ApiPostHandlers.RegisterDefaultAPI("test3", func(args *handler.APIArgs) (ret interface{}, err error) {
		d, err := h.db.AnyCollection.GenerateRawStruct()
		log.Debug(d)
		return d, err
	})
	
	var userController = user.InitController(h.db.WXUser)
	h.DefaultApiHandler.InjectController(userController)
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

type AdminConfig struct {
	NeedAdmin      bool `json:"_admin"`
	NeedSuperAdmin bool `json:"_superAdmin"`
}

type AdminStructConfig struct {
	Admin      bool `json:"admin"`
	SuperAdmin bool `json:"superAdmin"`
}

type AdminPermissionConfig struct {
}

func (c AdminPermissionConfig) GetTableConfig() (reflect.Type) {
	return reflect.TypeOf(AdminConfig{})
}

func (c AdminPermissionConfig) GetFieldConfig() (reflect.Type) {
	return reflect.TypeOf(AdminStructConfig{})
}

func (h *ApiHandler) GetPermissionConfig() permission.PermissionConfig {
	return AdminPermissionConfig{}
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
