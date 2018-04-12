package handler

import (
	"github.com/bitly/go-simplejson"
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
	AuthAdminUser(string) bool
	AuthSuperAdminUser(string) (bool, bool)
}

type DefaultApiHandler struct {
	apiHandlers     ExtendApiHandler
	router          *gin.Engine
	ApiGetHandlers  JsonAPIFuncRoute
	ApiPostHandlers JsonAPIFuncRoute
	db              DB
	Config          ApiHandlerConfig
	TableConfig     permission.TableConfigMap
}

func NewDefaultHandlerFromConfig(config ApiHandlerConfig, ah ExtendApiHandler) {
	h := &DefaultApiHandler{
		Config:          config,
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

func (h *DefaultApiHandler) SetDefaultApiHandlerAndMountConfig() {
	vi := reflect.ValueOf(h.apiHandlers).Elem()
	tc := reflect.ValueOf(h.Config)
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
	var function JsonAPIFunc
	var exists bool
	userSession := h.GetSession(c)
	apiName := c.Param("api")
	function, exists = h.GetApiFunc(c.Request.Method, apiName)
	if exists == true {
		var jsonData = new(jsonx.Json)
		if c.Request.Method == http.MethodPost {
			jsonData.Json, err = simplejson.NewFromReader(c.Request.Body)
		} else {
			jsonData.Json = simplejson.New()
		}
		if err == nil {
			ret, err = function(c, jsonData, userSession)
			if err == nil {
				ok = true
			}
		}
		
	}
	if h.CheckDataBaseConnection(err); err == nil {
		ret = h.RenderPermission(c, userSession, ret)
	} else {
		log.Errorf("JsonAPI(%s) err = %v", apiName, err)
	}
	RenderJson(c, ok, ret, err)
}

type JsonAPIFuncRoute map[string]JsonAPIFunc

func NewJsonAPIFuncRoute() JsonAPIFuncRoute {
	return make(JsonAPIFuncRoute)
}

func (j JsonAPIFuncRoute) RegisterAPI(name string, function interface{}) {
	switch function := function.(type) {
	case JsonAPIFunc:
		j.RegisterJsonAPI(name, function)
	case func(*gin.Context, *jsonx.Json, *session.UserSession) (interface{}, error):
		j.RegisterJsonAPI(name, function)
	case DefaultAPI:
		j.RegisterDefaultAPI(name, function)
	case func(args APIArgs) (ret interface{}, err error):
		j.RegisterDefaultAPI(name, function)
	default:
		panic(function)
	}
}

func (j JsonAPIFuncRoute) RegisterJsonAPI(name string, function JsonAPIFunc) {
	j[name] = function
}

func (j JsonAPIFuncRoute) RegisterDefaultAPI(name string, api DefaultAPI) {
	j.RegisterAPI(name, func(context *gin.Context, json *jsonx.Json, userSession *session.UserSession) (i interface{}, e error) {
		return api(APIArgs{
			context, json, userSession,
		})
	})
}

func RenderJson(c *gin.Context, ok bool, data interface{}, err error) {
	c.Header("Access-Control-Allow-Origin", "*")
	if ok == true {
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

type JsonAPIFunc func(g *gin.Context, j *jsonx.Json, s *session.UserSession) (interface{}, error)

func (h *DefaultApiHandler) GetApiFunc(method, apiName string) (function JsonAPIFunc, exists bool) {
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
