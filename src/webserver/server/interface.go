package server

import (
	"github.com/gin-gonic/gin"
	"net"
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
	Addr         string
	server       *http.Server
	listener     *net.TCPListener
	ApiHandlers  ApiHandlers
	ServerConfig ServerConfig
}
