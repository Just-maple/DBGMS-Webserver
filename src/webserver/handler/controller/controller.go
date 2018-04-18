package controller

import (
	"webserver/handler"
	"net/http"
)

type DefaultController struct {
	handler handler.HandlerGetter
}

func (c *DefaultController) GetDefaultController() (handler.HandlerController) {
	return c
}

func (c *DefaultController) Init() {
}

func (c *DefaultController) InjectHandler(h handler.HandlerGetter) {
	c.handler = h
}

func (c *DefaultController) MakeRegisterGroupByMethod(method string, pm ...handler.PermissionAuth) *handler.RegisterGroup {
	return c.handler.GetApiHandlersFromMethod(method).MakeRegisterGroup(pm...)
}

func (c *DefaultController) RegisterApi(method, api string, function handler.DefaultAPIFunc, pm ...handler.PermissionAuth) {
	h := c.handler.GetApiHandlersFromMethod(method)
	h.RegisterDefaultAPI(api, function, pm...)
}

func (c *DefaultController) RegisterGetApi(addr string, function func(args *handler.APIArgs) (ret interface{}, err error), pm ...handler.PermissionAuth) {
	c.RegisterApi(http.MethodGet, addr, function, pm...)
	return
}
func (c *DefaultController) RegisterPostApi(addr string, function func(args *handler.APIArgs) (ret interface{}, err error), pm ...handler.PermissionAuth) {
	c.RegisterApi(http.MethodPost, addr, function, pm...)
	return
}
