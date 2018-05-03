package handler

import (
	"github.com/gin-gonic/gin"
	"logger"
	"reflect"
	. "webserver/args"
	"webserver/dbx"
	"webserver/errorx"
	"webserver/permission"
	ws "webserver/server"
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
	GetPermissionConfig() *permission.PermissionConfig
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
	TableController   TableController
}

type TableController interface {
	HandlerController
	GetPermissionConfig() *permission.Config
	SetAccessConfig(func(args *APIArgs) permission.AccessConfig)
	GetConfigTableFromMString(args *APIArgs)(tableConfig map[string]string, err error)
	SetPath(string)
}

func newJsonAPIFuncRoute() JsonAPIFuncRoute {
	return make(JsonAPIFuncRoute)
}

func NewDefaultHandlerFromConfig(config ApiHandlerConfig, ah ExtendApiHandler) {
	h := &DefaultApiHandler{
		config:            config,
		ApiGetHandlers:    newJsonAPIFuncRoute(),
		ApiPostHandlers:   newJsonAPIFuncRoute(),
		ApiPutHandlers:    newJsonAPIFuncRoute(),
		ApiDeleteHandlers: newJsonAPIFuncRoute(),
		apiHandlers:       ah,
	}
	h.setDefaultApiHandlerAndMountConfig()
	return
}

func (h *DefaultApiHandler) InitMetaConfig() {
}

func (h *DefaultApiHandler) setDefaultApiHandlerAndMountConfig() {
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
	//err = h.InjectTableController(h.apiHandlers.GetPermissionConfig())
	//if err != nil {
	//	log.Fatal("Init TableConfig From Files Error = ", err)
	//}
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
