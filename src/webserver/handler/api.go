package handler

import (
	"github.com/gin-gonic/gin"
	"time"
	"webserver/jsonx"
	"webserver/session"
	"webserver/utilsx"
	"net/http"
	"github.com/bitly/go-simplejson"
)

type JsonAPIFunc func(g *gin.Context, j *jsonx.Json, s *session.UserSession) (interface{}, error)

type DefaultAPIFunc func(args *APIArgs) (ret interface{}, err error)
type PermissionAuth func(*session.UserSession) (bool)

type DefaultAPI struct {
	DefaultAPIFunc
	PermissionAuth []PermissionAuth
}

type APIArgs struct {
	context *gin.Context
	json    *jsonx.Json
	session *session.UserSession
}

func (api *DefaultAPI) Run(c *gin.Context, userSession *session.UserSession) (ret interface{}, err error) {
	var jsonData = new(jsonx.Json)
	if c.Request.Method == http.MethodPost {
		jsonData.Json, err = simplejson.NewFromReader(c.Request.Body)
	} else {
		jsonData.Json = simplejson.New()
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

type RegisterGroup struct {
	Route *JsonAPIFuncRoute
	pm    *[]PermissionAuth
}

func (r *RegisterGroup) RegisterDefaultAPI(name string, function DefaultAPIFunc) {
	r.Route.RegisterDefaultAPI(name, function, *r.pm...)
}
func (r *RegisterGroup) RegisterAPI(name string, function JsonAPIFunc) {
	r.Route.RegisterAPI(name, function, *r.pm...)
}

func (j *JsonAPIFuncRoute) MakeRegisterGroup(pm ...PermissionAuth) *RegisterGroup {
	return &RegisterGroup{j, &pm}
}

func (j *JsonAPIFuncRoute) RegisterAPI(name string, function JsonAPIFunc, pm ...PermissionAuth) {
	j.registerJsonAPI(name,
		func(args *APIArgs) (i interface{}, e error) {
			return function(args.context, args.json, args.session)
		}, pm)
}

func (j JsonAPIFuncRoute) registerJsonAPI(name string, function DefaultAPIFunc, pm []PermissionAuth) {
	if j[name] != nil {
		panic("route already existed")
	} else {
		j[name] = &DefaultAPI{
			function, pm,
		}
	}
}

func (j *JsonAPIFuncRoute) RegisterDefaultAPI(name string, api DefaultAPIFunc, pm ...PermissionAuth) {
	j.registerJsonAPI(name, api, pm)
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
