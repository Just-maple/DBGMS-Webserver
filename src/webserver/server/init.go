package server

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"time"
	"webserver/logger"
)

var log = logger.Log

func NewWebServer(config ServerConfig) (svr *WebServer) {
	svr = &WebServer{
		ServerConfig: config,
		Addr:         config.GetServerAddr(),
		server:       &http.Server{ReadHeaderTimeout: time.Second * 30, WriteTimeout: time.Second * 30},
	}
	return
}

func (svr *WebServer) Start() (err error) {
	svr.server.SetKeepAlivesEnabled(true)
	addr, err := net.ResolveTCPAddr("tcp4", svr.Addr)
	if err != nil {
		log.Fatalf("WebServer::Start net.ResolveTCPAddr err(%v)", err)
		return
	}
	listener, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		log.Fatalf("WebServer::Start net.ListenTCP err(%v)", err)
		return
	}
	svr.server.Serve(listener)
	return
}

func (svr *WebServer) InitSession() gin.HandlerFunc {
	store := sessions.NewCookieStore([]byte(svr.ServerConfig.GetSessionSecretKey()))
	store.Options(sessions.Options{MaxAge: 60 * 60 * 24 * 2})
	return sessions.Sessions(svr.ServerConfig.GetSessionKey(), store)
}

func (svr *WebServer) InitRouter() *gin.Engine {
	r := gin.Default()
	svr.server.Handler = r
	svr.ApiHandlers.SetRouter(r)
	return r
}

func (svr *WebServer) InitHandler(handler ApiHandlers) {
	svr.ApiHandlers = handler
	r := svr.InitRouter()
	r.Use(svr.InitSession())
	svr.ApiHandlers.InitDataBase()
	svr.ApiHandlers.RegisterDefaultAPI(svr.JsonAPI)
	svr.ApiHandlers.InitMetaConfig()
}
