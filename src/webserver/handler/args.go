package handler

import (
	"github.com/gin-gonic/gin"
	"webserver/jsonx"
	"webserver/session"
	"time"
	"webserver/utilsx"
	"gopkg.in/mgo.v2/bson"
)

type APIArgs struct {
	context *gin.Context
	Json    *jsonx.Json
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
	return arg.Json.Unmarshal(in)
}

func (arg *APIArgs) JsonKeyId() string {
	return arg.Json.GetStringId()
}

func (arg *APIArgs) SetUserId(userId bson.ObjectId) {
	arg.session.SetUserId(userId.Hex())
}

func (arg *APIArgs) JsonKey(key string) *jsonx.Json {
	return arg.Json.Get(key)
}