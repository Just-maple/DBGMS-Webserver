package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ServerConfig interface {
	GetServerAddr() string
	GetSessionKey() string
	GetSessionSecretKey() string
}

type ApiHandlers interface {
	RegisterJsonAPI()
	InitDataBase()
	InitMetaConfig()
	SetRouter(r *gin.Engine)
}

type WebServer struct {
	addr         string
	server       *http.Server
	apiHandlers  ApiHandlers
	serverConfig ServerConfig
}
