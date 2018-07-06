package controller

import (
	"net/http"
	"webserver/args"
	"webserver/dbx"
	"webserver/handler"
)

type DefaultController struct {
	handler    handler.HandlerGetter
	Collection *dbx.Collection
}

func NewDefaultController(collection *dbx.Collection) *DefaultController {
	return &DefaultController{
		Collection: collection,
	}
}

func (c *DefaultController) Init() {
}

func (c *DefaultController) InjectHandler(h handler.HandlerGetter) {
	c.handler = h
}

func (c *DefaultController) MakeRegisterGroupByMethod(method string, pm ...handler.PermissionAuth) *handler.RegisterGroup {
	h := c.handler.GetApiHandlersFromMethod(method)
	return h.MakeRegisterGroup(pm...)
}

func (c *DefaultController) RegisterApi(method, api string, function handler.DefaultAPIFunc, pm ...handler.PermissionAuth) {
	h := c.handler.GetApiHandlersFromMethod(method)
	h.RegisterDefaultAPI(api, function, pm...)
}

func (c *DefaultController) RegisterGetApi(addr string, function func(args *args.APIArgs) (ret interface{}, err error), pm ...handler.PermissionAuth) {
	c.RegisterApi(http.MethodGet, addr, function, pm...)
	return
}

func (c *DefaultController) RegisterPostApi(addr string, function func(args *args.APIArgs) (ret interface{}, err error), pm ...handler.PermissionAuth) {
	c.RegisterApi(http.MethodPost, addr, function, pm...)
	return
}
