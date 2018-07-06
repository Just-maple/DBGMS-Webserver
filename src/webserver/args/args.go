package args

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
	"webserver/args/jsonx"
	"webserver/args/session"
	"webserver/permission"
	"webserver/utilsx"
)

type APIArgs struct {
	context *gin.Context
	Json    *jsonx.Json
	session *session.UserSession
}

func New(c *gin.Context, j *jsonx.Json, s *session.UserSession) *APIArgs {
	return &APIArgs{c, j, s}
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

func (arg *APIArgs) DefaultAPI() (*gin.Context, *jsonx.Json, *session.UserSession) {
	return arg.context, arg.Json, arg.session
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

func (arg *APIArgs) TransAjaxQuery() (matcherMap map[string]interface{}, keys []string, skipCnt, limitCnt int, sortKey, reverse string, tSTime, tETime time.Time, err error) {
	keys = arg.JsonKey("keys").MustStringArray()
	matcherMap = arg.JsonKey("matcher").MustMap()

	skip := arg.context.DefaultQuery("skip", "0")
	skipCnt, err = strconv.Atoi(skip)
	if err != nil {
		skipCnt = 0
		err = nil
	}
	limit := arg.context.DefaultQuery("limit", "10")
	limitCnt, err = strconv.Atoi(limit)
	if err != nil {
		limitCnt = 0
		err = nil
	}
	sortKey = arg.context.DefaultQuery("sort", "")
	reverse = arg.context.DefaultQuery("reverse", "")
	tSTime, tETime = arg.Time()
	return
}
