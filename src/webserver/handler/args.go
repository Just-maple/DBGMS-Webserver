package handler

import (
	"github.com/gin-gonic/gin"
	"webserver/jsonx"
	"webserver/session"
	"time"
	"webserver/utilsx"
)

type APIArgs struct {
	context *gin.Context
	json    *jsonx.Json
	session *session.UserSession
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

func (arg *APIArgs) JsonUnmarshal(in interface{}) error {
	return arg.json.Unmarshal(in)
}

func (arg *APIArgs) JsonKeyId() string {
	return arg.json.GetStringId()
}

func (arg *APIArgs) JsonKey(key string) *jsonx.Json {
	return arg.json.Get(key)
}

