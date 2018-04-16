package handler

import (
	"github.com/bitly/go-simplejson"
	"github.com/gin-gonic/gin"
	"time"
	"webserver/jsonx"
	"webserver/session"
	"webserver/utilsx"
)

type JsonAPIFunc func(g *gin.Context, j *jsonx.Json, s *session.UserSession) (interface{}, error)

type DefaultAPI func(args *APIArgs) (ret interface{}, err error)
type APIArgs struct {
	context *gin.Context
	json    *jsonx.Json
	session *session.UserSession
}

func (j *JsonAPIFuncRoute) RegisterAPI(name string, function JsonAPIFunc) {
	j.registerJsonAPI(name,
		func(args *APIArgs) (i interface{}, e error) {
			return function(args.context, args.json, args.session)
		})
}

func (j JsonAPIFuncRoute) registerJsonAPI(name string, function DefaultAPI) {
	if j[name] != nil {
		panic("route already existed")
	} else {
		j[name] = function
	}
}

func (j *JsonAPIFuncRoute) RegisterDefaultAPI(name string, api DefaultAPI) {
	j.registerJsonAPI(name, api)
}

func (arg *APIArgs) Time() (st, et time.Time) {
	return utilsx.TransTime(arg.context)
}

func (arg *APIArgs) IP() string {
	return arg.context.ClientIP()
}

func (arg *APIArgs) Query(key string) string {
	return arg.context.Query(key)
}

func (arg *APIArgs) UserId() (valid bool, userId string) {
	return arg.session.AuthUserSession()
}

func (arg *APIArgs) JsonKeyId() string {
	return arg.json.GetStringId()
}

func (arg *APIArgs) JsonKey(key string) *simplejson.Json {
	return arg.json.Get(key)
}
