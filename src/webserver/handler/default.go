package handler

import (
	"github.com/gin-gonic/gin"
	"reflect"
	"webserver/dbx"
	"webserver/errorx"
	"logger"
	"webserver/permission"
	ws "webserver/server"
	. "webserver/args"
)

var log = logger.Log
var ErrAuthFailed = errorx.ErrAuthFailed

type ApiHandlerConfig interface {
	GetTablePath() string
	GetMgoDBUrl() string
}

type ExtendApiHandler interface {
	ws.ApiHandlers
	RegisterAPI()
	NewDataBase() DB
	GetPermissionConfig() permission.PermissionConfig
	GetAccessConfig(args *APIArgs) permission.AccessConfig
}

type DB interface {
}

var _ ws.ApiHandlers = &DefaultApiHandler{}

type DefaultApiHandler struct {
	apiHandlers       ExtendApiHandler
	router            *gin.Engine
	ApiGetHandlers    JsonAPIFuncRoute
	ApiPostHandlers   JsonAPIFuncRoute
	ApiPutHandlers    JsonAPIFuncRoute
	ApiDeleteHandlers JsonAPIFuncRoute
	db                DB
	config            ApiHandlerConfig
	TableController   *TableController
}

func NewJsonAPIFuncRoute() JsonAPIFuncRoute {
	return make(JsonAPIFuncRoute)
}

func NewDefaultHandlerFromConfig(config ApiHandlerConfig, ah ExtendApiHandler) {
	h := &DefaultApiHandler{
		config:            config,
		ApiGetHandlers:    NewJsonAPIFuncRoute(),
		ApiPostHandlers:   NewJsonAPIFuncRoute(),
		ApiPutHandlers:    NewJsonAPIFuncRoute(),
		ApiDeleteHandlers: NewJsonAPIFuncRoute(),
		apiHandlers:       ah,
	}
	h.SetDefaultApiHandlerAndMountConfig()
	var err error
	h.TableController, err = InjectTableController(h, ah.GetPermissionConfig())
	if err != nil {
		log.Fatal("Init TableConfig FromFiles Error = ", err)
	}
	return
}

func (h *DefaultApiHandler) InitMetaConfig() {
}

func (h *DefaultApiHandler) SetDefaultApiHandlerAndMountConfig() {
	vi := reflect.ValueOf(h.apiHandlers).Elem()
	tc := reflect.ValueOf(h.config)
	fi := vi.NumField()
	ht := reflect.TypeOf(h)
	var flag bool
	for i := 0; i < fi; i++ {
		if !vi.Field(i).CanSet() {
			continue
		}
		switch vi.Field(i).Type() {
		case tc.Type():
			vi.Field(i).Set(tc)
		case ht:
			if !flag {
				vi.Field(i).Set(reflect.ValueOf(h))
				flag = true
			} else {
				panic("Found More than One Default ApiHandler")
			}
		}
		
	}
	if !flag {
		panic("Not Found Default ApiHandler")
	}
	return
}

func (h *DefaultApiHandler) InitDataBase() {
	var err error
	db := h.apiHandlers.NewDataBase()
	h.db = db
	err = dbx.NewMgoDB(h.config.GetMgoDBUrl(), db)
	if err != nil {
		log.Fatal("Init MgoDataBase Error = ", err)
		return
	}
	return
}

func (h *DefaultApiHandler) IsDataBaseConnectionError(err error) bool {
	return err != nil && (err.Error() == "Closed explicitly" || err.Error() == "EOF")
}

func (h *DefaultApiHandler) CheckDataBaseConnection(err error) {
	if h.IsDataBaseConnectionError(err) {
		h.apiHandlers.InitDataBase()
	}
}
