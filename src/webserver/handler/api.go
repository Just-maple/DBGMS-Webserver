package handler

import (
	"github.com/gin-gonic/gin"
	"webserver/jsonx"
	"webserver/session"
	"net/http"
)

type JsonAPIFunc func(g *gin.Context, j *jsonx.Json, s *session.UserSession) (interface{}, error)

type DefaultAPIFunc func(args *APIArgs) (ret interface{}, err error)
type PermissionAuth func(*session.UserSession) (bool)

type DefaultAPI struct {
	DefaultAPIFunc
	PermissionAuth []PermissionAuth
}

func (api *DefaultAPI) Run(c *gin.Context, userSession *session.UserSession) (ret interface{}, err error) {
	var jsonData *jsonx.Json
	if c.Request.Method == http.MethodPost {
		jsonData, err = jsonx.NewFromReader(c.Request.Body)
	} else {
		jsonData = jsonx.New()
	}
	if err == nil {
		if len(api.PermissionAuth) > 0 {
			for i := range api.PermissionAuth {
				valid := api.PermissionAuth[i](userSession)
				if !valid {
					err = ErrAuthFailed
					return
				}
			}
		}
		ret, err = api.DefaultAPIFunc(&APIArgs{c, jsonData, userSession})
	}
	return
}
