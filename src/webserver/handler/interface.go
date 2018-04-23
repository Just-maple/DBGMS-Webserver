package handler

import (
	. "webserver/args"
	"webserver/permission"
)

type HandlerGetter interface {
	GetApiHandlersFromMethod(method string) (handler JsonAPIFuncRoute)
	InjectController(c HandlerController)
}

type TableHandler interface {
	GetApiHandlersFromMethod(method string) (handler JsonAPIFuncRoute)
	GetAccessConfigFromArgs(args *APIArgs) (access permission.AccessConfig)
	GetTablePath() string
	GetAccessConfig(args *APIArgs) permission.AccessConfig
}

type HandlerController interface {
	GetDefaultController() HandlerController
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
	c.GetDefaultController().InjectHandler(h)
	c.Init()
}
