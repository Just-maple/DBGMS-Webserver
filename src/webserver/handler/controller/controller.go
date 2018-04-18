package controller

import "webserver/handler"

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

func (c *DefaultController) RegisterApi(method, api string, function handler.DefaultAPIFunc, pm ...handler.PermissionAuth) {
	h := c.handler.GetApiHandlersFromMethod(method)
	h.RegisterDefaultAPI(api, function, pm...)
}
