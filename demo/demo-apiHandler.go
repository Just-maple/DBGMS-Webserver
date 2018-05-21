package main

import (
	"logger"
	"webserver/args"
	"webserver/handler"
	"webserver/controller/user"
	"webserver/controller/table"
	"webserver/dbx"
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
	UserController  *user.Controller
	TableController *table.Controller
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
	
	//init default user controller with your user database collection
	h.UserController = user.InitController(h.db.WXUser)
	h.InjectController(h.UserController)
	
	//inject your table controller by collection store
	h.TableController = table.InitAdminTableControllerByCollection(h.db.PermissionTable)
	// or user path to store your table config
	//h.TableController = table.InitAdminTableControllerByPath(h.Config.GetTablePath())
	h.InjectTableController(h.TableController)
	//you can handle your extend config here
	//this method will execute after database init
	//and before server start
}

type test struct {
	Yttt string `json:"yttt" bson:"testb"`
	Hhh  int `bson:"fff"`
}

func (h *ApiHandler) NewDataBase() dbx.DB {
	//implement method NewDataBase
	log.Debug("New DataBase From Handler")
	//init your DataBase
	h.db = new(DataBase)
	h.db.RegisterStruct("ttt", test{})
	return h.db
	//should return interface implement server.DB
}
