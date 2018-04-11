package server

import (
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"webserver/session"
)

type ServerConfig interface {
	GetServerAddr() string
	GetSessionKey() string
	GetSessionSecretKey() string
}

type DB interface {
	AuthAdminUser(string) bool
	AuthSuperAdminUser(string) (bool, bool)
}

type ApiHandlers interface {
	GetApiFunc(string, string) (JsonAPIFunc, bool)

	RegisterDefaultAPI(gin.HandlerFunc)
	RegisterAPI()

	InitDataBase()
	NewDataBase() DB
	InitMetaConfig()
	CheckDataBaseConnection(err error)
	RenderPermission(c *gin.Context, session *session.UserSession, in interface{}) (out interface{})
	SetRouter(r *gin.Engine)
	GetSession(c *gin.Context) (us *session.UserSession)
}

type WebServer struct {
	Addr         string
	server       *http.Server
	listener     *net.TCPListener
	ApiHandlers  ApiHandlers
	ServerConfig ServerConfig
}
