package handler

import (
	"github.com/gin-gonic/gin"
	. "webserver/args"
	"webserver/args/jsonx"
	"webserver/args/session"
)

type JsonAPIFunc func(g *gin.Context, j *jsonx.Json, s *session.UserSession) (interface{}, error)

type DefaultAPIFunc func(args *APIArgs) (ret interface{}, err error)
type PermissionAuth func(*APIArgs) bool

type DefaultAPI struct {
	DefaultAPIFunc
	PermissionAuth []PermissionAuth
}

func (api *DefaultAPI) run(args *APIArgs) (ret interface{}, err error) {
	for i := range api.PermissionAuth {
		if !api.PermissionAuth[i](args) {
			return false, ErrAuthFailed
		}
	}
	ret, err = api.DefaultAPIFunc(args)
	return
}
