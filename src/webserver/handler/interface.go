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

type TableController interface {
	HandlerController
	GetPermissionConfig() *permission.Config
	SetAccessConfig(func(args *APIArgs) permission.AccessConfig)
	GetConfigTableFromArgs(args *APIArgs) (tableConfig map[string]string, err error)
}

func (h *DefaultApiHandler) GetAccessConfig(args *APIArgs) permission.AccessConfig {
	return h.apiHandlers.GetAccessConfig(args)
}

func (h *DefaultApiHandler) InjectController(c HandlerController) {
	c.InjectHandler(h)
	c.Init()
}

func (h *DefaultApiHandler) InjectTableController(c TableController) {
	h.TableController = c
	c.SetAccessConfig(h.GetAccessConfig)
	c.InjectHandler(h)
	c.Init()
}
