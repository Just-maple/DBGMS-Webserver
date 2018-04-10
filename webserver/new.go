package webserver

import (
	. "webserver/handler"
	. "webserver/server"
)

type Config interface {
	ApiHandlerConfig
	ServerConfig
}

func NewWebServerFromHandler(config Config, apiHandler ApiHandlers) (svr *WebServer) {
	NewDefaultHandlerFromConfig(config, apiHandler)
	svr = NewWebServer(config)
	svr.InitHandler(apiHandler)
	return
}
