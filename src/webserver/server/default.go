package server

import (
	"github.com/bitly/go-simplejson"
	"github.com/gin-gonic/gin"
	"webserver/jsonx"
	"webserver/session"
)

type DefaultAPI func(args DefaultAPIArgs) (ret interface{}, err error)

type DefaultAPIArgs struct {
	context *gin.Context
	json    *jsonx.Json
	session *session.UserSession
}

func (h JsonAPIFuncRoute) RegisterDefaultAPI(name string, api DefaultAPI) {
	h.RegisterAPI(name, func(context *gin.Context, json *jsonx.Json, userSession *session.UserSession) (i interface{}, e error) {
		return api(DefaultAPIArgs{
			context, json, userSession,
		})
	})
}

func (arg *DefaultAPIArgs) GetQuery(key string) string {
	return arg.context.Query(key)
}

func (arg *DefaultAPIArgs) GetUserId() (valid bool, userId string) {
	return arg.session.AuthUserSession()
}

func (arg *DefaultAPIArgs) GetJsonKeyId() string {
	return arg.json.GetStringId()
}

func (arg *DefaultAPIArgs) GetJsonKey(key string) *simplejson.Json {
	return arg.json.Get(key)
}
