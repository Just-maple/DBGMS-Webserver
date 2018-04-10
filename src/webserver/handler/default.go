package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"webserver/dbx"
	"webserver/errorx"
	"webserver/jsonx"
	"webserver/logger"
	"webserver/permission"
	ws "webserver/server"
	"webserver/session"
)

var log = logger.Log
var ErrAuthFailed = errorx.ErrAuthFailed

type ApiHandlerConfig interface {
	GetTablePath() string
	GetMgoDBUrl() string
}

type DefaultApiHandler struct {
	apiHandlers     ws.ApiHandlers
	router          *gin.Engine
	ApiGetHandlers  ws.JsonAPIFuncRoute
	ApiPostHandlers ws.JsonAPIFuncRoute
	db              ws.DB
	Config          ApiHandlerConfig
	TableConfig     permission.TableConfigMap
}

func NewDefaultHandlerFromConfig(config ApiHandlerConfig, ah ws.ApiHandlers) {
	h := &DefaultApiHandler{
		Config:          config,
		ApiGetHandlers:  ws.NewJsonAPIFuncRoute(),
		ApiPostHandlers: ws.NewJsonAPIFuncRoute(),
		apiHandlers:     ah,
	}
	ptrC := reflect.ValueOf(ah).Elem().FieldByName("Config")
	if ptrC.IsValid() && ptrC.CanSet() {
		ptrC.Set(reflect.ValueOf(config))
	}
	ptrH := reflect.ValueOf(ah).Elem().FieldByName("DefaultApiHandler")
	if ptrH.IsValid() && ptrH.CanSet() && ptrH.Type() == reflect.TypeOf(h) {
		ptrH.Set(reflect.ValueOf(h))
	} else {
		panic("Invalid ApiHandler")
	}
	err := h.InitAllConfigTableFromFiles()
	if err != nil {
		log.Fatal("Init TableConfig FromFiles Error = ", err)
	}
	return
}

func (h *DefaultApiHandler) RegisterRouter(method, path string, function gin.HandlerFunc) {
	switch method {
	case http.MethodGet:
		h.router.GET(path, function)
	case http.MethodPost:
		h.router.POST(path, function)
	}
	
}

func (h *DefaultApiHandler) SetRouter(r *gin.Engine) {
	h.router = r
}

func (h *DefaultApiHandler) RegisterAPI(api gin.HandlerFunc) {
	h.router.GET("/api/:api", api)
	h.router.POST("/api/:api", api)
	h.ApiPostHandlers.RegisterAPI("test", h.Test)
	h.ApiGetHandlers.RegisterAPI("test", h.Test)
	h.apiHandlers.RegisterGETAPI()
	h.apiHandlers.RegisterPOSTAPI()
	h.apiHandlers.RegisterSpecificAPI()
	//for _,v:=range h.ApiGetHandlers{
	//	log.Debug(strings.Split(runtime.FuncForPC(reflect.ValueOf(v).Pointer()).Name(),"."))
	//}
}

func (h *DefaultApiHandler) GetApiFunc(method, apiName string) (function ws.JsonAPIFunc, exists bool) {
	switch method {
	case http.MethodGet:
		function, exists = h.ApiGetHandlers[apiName]
	case http.MethodPost:
		function, exists = h.ApiPostHandlers[apiName]
	}
	return
}

func (h *DefaultApiHandler) InitDataBase() {
	var err error
	db := h.apiHandlers.NewDataBase()
	h.db = db
	err = dbx.NewMgoDB(h.Config.GetMgoDBUrl(), db)
	if err != nil {
		log.Fatal("Init MgoDataBase Error = ", err)
		return
	}
	return
}

func (h *DefaultApiHandler) Test(c *gin.Context, j *jsonx.Json, us *session.UserSession) (ret interface{}, err error) {
	ret = j.Get("test").MustString()
	if ret == "" {
		ret = "test success"
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
