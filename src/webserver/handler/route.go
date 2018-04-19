package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"webserver/jsonx"
	"webserver/session"
)

type JsonAPIFuncRoute map[string]*DefaultAPI

func (h *DefaultApiHandler) RegisterRouter(method, path string, function gin.HandlerFunc) {
	h.router.Handle(method, path, function)
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

func NewJsonAPIFuncRoute() JsonAPIFuncRoute {
	return make(JsonAPIFuncRoute)
}

func (h *DefaultApiHandler) GetApiHandlersFromMethod(method string) (handler JsonAPIFuncRoute) {
	switch method {
	case http.MethodGet:
		return h.ApiGetHandlers
	case http.MethodPost:
		return h.ApiPostHandlers
	case http.MethodPut:
		return h.ApiPutHandlers
	case http.MethodDelete:
		return h.ApiDeleteHandlers
	default:
		panic("method invalid " + method)
		return
	}
}
func (h *DefaultApiHandler) getApiFunc(method, apiName string) (function *DefaultAPI, exists bool) {
	function, exists = h.GetApiHandlersFromMethod(method)[apiName]
	return
}
func (h *DefaultApiHandler) Test(c *gin.Context, j *jsonx.Json, us *session.UserSession) (ret interface{}, err error) {
	ret = j.Get("test").MustString()
	if ret == "" {
		ret = "test success"
	}
	return
}
