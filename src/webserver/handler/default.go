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
	GetMgoDBUrl() string
}

type ExtendApiHandler interface {
	ws.ApiHandlers
	RegisterAPI()
	NewDataBase() dbx.DB
	GetAccessConfig(args *APIArgs) permission.AccessConfig
}

var _ ws.ApiHandlers = &DefaultApiHandler{}

type DefaultApiHandler struct {
	apiHandlers       ExtendApiHandler
	router            *gin.Engine
	ApiGetHandlers    JsonAPIFuncRoute
	ApiPostHandlers   JsonAPIFuncRoute
	ApiPutHandlers    JsonAPIFuncRoute
	ApiDeleteHandlers JsonAPIFuncRoute
	db                dbx.DB
	config            ApiHandlerConfig
	TableController   TableController
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
