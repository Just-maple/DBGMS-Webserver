package handler

import (
	"github.com/gin-gonic/gin"
	"webserver/jsonx"
	"webserver/session"
	. "webserver/args"
)

type JsonAPIFunc func(g *gin.Context, j *jsonx.Json, s *session.UserSession) (interface{}, error)

type DefaultAPIFunc func(args *APIArgs) (ret interface{}, err error)
type PermissionAuth func(*APIArgs) bool

type DefaultAPI struct {
	DefaultAPIFunc
	PermissionAuth []PermissionAuth
}

func (api *DefaultAPI) Run(args *APIArgs) (ret interface{}, err error) {
	for i := range api.PermissionAuth {
		if !api.PermissionAuth[i](args) {
			return false, ErrAuthFailed
		}
	}
	ret, err = api.DefaultAPIFunc(args)
	return
}
