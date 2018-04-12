package handler

import (
	"github.com/bitly/go-simplejson"
	"github.com/gin-gonic/gin"
	"time"
	"webserver/jsonx"
	"webserver/session"
	"webserver/utilsx"
)

type DefaultAPI func(args APIArgs) (ret interface{}, err error)

type APIArgs struct {
	context *gin.Context
	json    *jsonx.Json
	session *session.UserSession
}

func (arg *APIArgs) Time() (st, et time.Time) {
	return utilsx.TransTime(arg.context)
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
