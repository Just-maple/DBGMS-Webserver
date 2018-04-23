package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"webserver/args/jsonx"
	"webserver/args/session"
)

type JsonAPIFuncRoute map[string]*DefaultAPI

func (h *DefaultApiHandler) RegisterRouter(method, path string, function gin.HandlerFunc) {
	h.router.Handle(method, path, function)
}

func (h *DefaultApiHandler) SetRouter(r *gin.Engine) {
	h.router = r
}

func (h *DefaultApiHandler) RegisterJsonAPI() {
	h.router.GET("/api/:api", h.jsonAPI)
	h.router.POST("/api/:api", h.jsonAPI)
	h.ApiPostHandlers.RegisterAPI("test", h.test)
	h.ApiGetHandlers.RegisterAPI("test", h.test)
	h.apiHandlers.RegisterAPI()
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
func (h *DefaultApiHandler) test(c *gin.Context, j *jsonx.Json, us *session.UserSession) (ret interface{}, err error) {
	ret = j.Get("test").MustString()
	if ret == "" {
		ret = "test success method:" + c.Request.Method
	}
	return
}
