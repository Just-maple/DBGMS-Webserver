package server

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"logger"
	"net"
	"net/http"
	"time"
)

var log = logger.Log

func NewWebServer(config ServerConfig) (svr *WebServer) {
	svr = &WebServer{
		serverConfig: config,
		addr:         config.GetServerAddr(),
		server:       &http.Server{ReadHeaderTimeout: time.Second * 30, WriteTimeout: time.Second * 30},
	}
	return
}

func (svr *WebServer) Start() (err error) {
	svr.server.SetKeepAlivesEnabled(true)
	addr, err := net.ResolveTCPAddr("tcp4", svr.addr)
	if err != nil {
		log.Fatalf("WebServer::Start net.ResolveTCPAddr err(%v)", err)
		return
	}
	listener, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		log.Fatalf("WebServer::Start net.ListenTCP err(%v)", err)
		return
	}
	log.Debugf("WebServer Start At Addr [ %v ]", svr.addr)
	svr.server.Serve(listener)
	return
}

func (svr *WebServer) initSession() gin.HandlerFunc {
	store := sessions.NewCookieStore([]byte(svr.serverConfig.GetSessionSecretKey()))
	store.Options(sessions.Options{MaxAge: 60 * 60 * 24 * 2})
	return sessions.Sessions(svr.serverConfig.GetSessionKey(), store)
}

func (svr *WebServer) initRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	svr.server.Handler = r
	svr.apiHandlers.SetRouter(r)
	return r
}

func (svr *WebServer) InitHandler(handler ApiHandlers) {
	svr.apiHandlers = handler
	r := svr.initRouter()
	r.Use(svr.initSession())
	svr.apiHandlers.InitDataBase()
	svr.apiHandlers.RegisterJsonAPI()
	svr.apiHandlers.InitMetaConfig()
}
