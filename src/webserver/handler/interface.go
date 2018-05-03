package handler

import (
	. "webserver/args"
	"webserver/permission"
)

type HandlerGetter interface {
	GetApiHandlersFromMethod(method string) (handler JsonAPIFuncRoute)
	InjectController(c HandlerController)
}

type HandlerController interface {
	InjectHandler(handler HandlerGetter)
	Init()
}

func (h *DefaultApiHandler) GetAccessConfig(args *APIArgs) permission.AccessConfig {
	return h.apiHandlers.GetAccessConfig(args)
}

func (h *DefaultApiHandler) GetTablePath() string {
	return h.config.GetTablePath()
}
func (h *DefaultApiHandler) InjectController(c HandlerController) {
	c.InjectHandler(h)
	c.Init()
}

func (h *DefaultApiHandler) InjectTableController(c TableController) {
	h.TableController = c
	c.SetAccessConfig(h.GetAccessConfig)
	c.InjectHandler(h)
	c.SetPath(h.GetTablePath())
	c.Init()
}
