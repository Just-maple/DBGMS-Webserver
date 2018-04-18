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

type ExtendApiHandler interface {
	ws.ApiHandlers
	RegisterAPI()
	NewDataBase() DB
}

type DB interface {
	GetAccessConfig(string) (permission.AccessConfig)
}

var _ ws.ApiHandlers = &DefaultApiHandler{}

type DefaultApiHandler struct {
	apiHandlers      ExtendApiHandler
	router           *gin.Engine
	ApiGetHandlers   JsonAPIFuncRoute
	ApiPostHandlers  JsonAPIFuncRoute
	db               DB
	config           ApiHandlerConfig
	PermissionConfig permission.TableConfigMap
}

func NewDefaultHandlerFromConfig(config ApiHandlerConfig, ah ExtendApiHandler) {
	h := &DefaultApiHandler{
		config:          config,
		ApiGetHandlers:  NewJsonAPIFuncRoute(),
		ApiPostHandlers: NewJsonAPIFuncRoute(),
		apiHandlers:     ah,
	}
	h.SetDefaultApiHandlerAndMountConfig()
	err := h.InitAllConfigTableFromFiles()
	if err != nil {
		log.Fatal("Init TableConfig FromFiles Error = ", err)
	}
	return
}

func (h *DefaultApiHandler) InitMetaConfig() {
	//if h.apiHandlers.InitMetaConfig != nil {
	//	h.apiHandlers.InitMetaConfig()
	//}
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

func (h *DefaultApiHandler) SetDefaultConfig(in interface{}) {
	vi := reflect.ValueOf(in).Elem()
	fi := vi.NumField()
	for i := 0; i < fi; i++ {
		switch vi.Field(i).Interface().(type) {
		case ApiHandlerConfig:
			vi.Field(i).Set(reflect.ValueOf(h))
		}
	}
	return
}

func (h *DefaultApiHandler) RegisterRouter(method, path string, function gin.HandlerFunc) {
	h.router.Handle(method, path, function)
}

func (h *DefaultApiHandler) SetRouter(r *gin.Engine) {
	h.router = r
}

func (h *DefaultApiHandler) RegisterAPI() {
	log.Info("you had not register any API,register your custom API by rewrite method RegisterAPI")
}

func (h *DefaultApiHandler) RegisterJsonAPI() {
	h.router.GET("/api/:api", h.JsonAPI)
	h.router.POST("/api/:api", h.JsonAPI)
	h.ApiPostHandlers.RegisterAPI("test", h.Test)
	h.ApiGetHandlers.RegisterAPI("test", h.Test)
	h.apiHandlers.RegisterAPI()
}

func (h *DefaultApiHandler) JsonAPI(c *gin.Context) {
	var ok bool
	var ret interface{}
	var err error
	var userSession = h.GetSession(c)
	var apiName = c.Param("api")
	function, exists := h.GetApiFunc(c.Request.Method, apiName)
	if exists {
		ret, err = function.Run(c, userSession)
		if h.CheckDataBaseConnection(err); err == nil {
			ret = h.RenderPermission(c, userSession, ret)
			ok = true
		} else {
			log.Errorf("JsonAPI(%s) err = %v", apiName, err)
		}
	}
	RenderJson(c, ok, ret, err)
}

type JsonAPIFuncRoute map[string]*DefaultAPI

func NewJsonAPIFuncRoute() JsonAPIFuncRoute {
	return make(JsonAPIFuncRoute)
}

func RenderJson(c *gin.Context, ok bool, data interface{}, err error) {
	c.Header("Access-Control-Allow-Origin", "*")
	if ok {
		c.JSON(http.StatusOK, map[string]interface{}{
			"ok":   ok,
			"data": data,
		})
	} else {
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, map[string]interface{}{
				"ok":   ok,
				"data": data,
				"err":  err.Error(),
			})
			
		} else {
			c.IndentedJSON(http.StatusInternalServerError, map[string]interface{}{
				"ok":   ok,
				"data": data,
			})
		}
	}
}

func (h *DefaultApiHandler) GetApiHandlersFromMethod(method string) (handler JsonAPIFuncRoute) {
	switch method {
	case http.MethodGet:
		return h.ApiGetHandlers
	case http.MethodPost:
		return h.ApiPostHandlers
	default:
		panic("method invalid " + method)
		return
	}
}
func (h *DefaultApiHandler) GetApiFunc(method, apiName string) (function *DefaultAPI, exists bool) {
	function, exists = h.GetApiHandlersFromMethod(method)[apiName]
	return
}

func (h *DefaultApiHandler) NewDataBase() DB {
	panic("you should provide a custom DB with AccessControl Config")
	return *new(DB)
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
