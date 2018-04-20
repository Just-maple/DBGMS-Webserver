package handler

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"time"
	"webserver/jsonx"
	"webserver/session"
	"webserver/utilsx"
	"webserver/permission"
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

func (arg *APIArgs) Api() string {
	return arg.context.Param("api")
}

func (arg *APIArgs) ClearSession() {
	arg.session.Clear()
	arg.session.Save()
}

func (arg *APIArgs) JsonKey(key string) *jsonx.Json {
	return arg.Json.Get(key)
}

func (arg *APIArgs) GetConfigTable(t *permission.Config) (config *permission.StructConfig, has bool) {
	table, has := t.TableMap[arg.Api()]
	if has {
		config = table.StructConfig
	}
	return
}
