package handler

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
	"webserver/dbx"
	"webserver/jsonx"
	. "webserver/server"
	"webserver/session"
	. "webserver/utilsx"
)

func TransAjaxQuery(c *gin.Context, j *jsonx.Json) (matcherMap map[string]interface{}, keys []string, skipCnt, limitCnt int, sortKey, reverse string, tSTime, tETime time.Time, err error) {
	keys = j.Get("keys").MustStringArray()
	matcherMap = j.Get("matcher").MustMap()

	skip := c.DefaultQuery("skip", "0")
	skipCnt, err = strconv.Atoi(skip)
	if err != nil {
		skipCnt = 0
		err = nil
	}
	limit := c.DefaultQuery("limit", "10")
	limitCnt, err = strconv.Atoi(limit)
	if err != nil {
		limitCnt = 0
		err = nil
	}
	sortKey = c.DefaultQuery("sort", "")
	reverse = c.DefaultQuery("reverse", "")
	tSTime, tETime = TransTime(c)
	return
}

func (h DefaultApiHandler) GetAjaxQuery(c *gin.Context, j *jsonx.Json) (res *dbx.AjaxQuery, err error) {
	matcherMap, keys, skipCnt, limitCnt, sortKey, reverse, tSTime, tETime, err := TransAjaxQuery(c, j)
	pmConfig, _ := h.TableConfig.GetConfigTableFromContext(c)
	res = &dbx.AjaxQuery{
		MatcherMap:       matcherMap,
		SortKey:          sortKey,
		SelectKeys:       keys,
		SkipCount:        skipCnt,
		LimitCount:       limitCnt,
		TimeStart:        tSTime,
		TimeEnd:          tETime,
		Reverse:          reverse,
		PermissionConfig: pmConfig,
	}
	return
}

func (h *DefaultApiHandler) GetDataByAjaxQuery(c *gin.Context, j *jsonx.Json, us *session.UserSession, ajaxConfig *dbx.AjaxStructConfig) (res map[string]interface{}, err error) {
	query, err := h.GetAjaxQuery(c, j)
	if err != nil {
		return
	}
	ia, is := h.AuthSuperAdminUserSession(us)
	data, count, err := query.AjaxSearch(ajaxConfig)
	if err != nil {
		log.Error(err)
		return
	}
	data = query.MakeAjaxReturnWithSelectKeysAndPermissionControl(data, ia, is)
	res = map[string]interface{}{
		"data": data,
		"cnt":  count,
	}
	return
}

func (h *DefaultApiHandler) RegisterAjaxJsonApi(dataApiAddr, distinctApiAddr string, configMaker func() dbx.AjaxStructConfig) {
	config := configMaker()
	h.ApiPostHandlers.RegisterAPI(dataApiAddr, h.GetAjaxApi(&config))
	h.ApiGetHandlers.RegisterAPI(distinctApiAddr, h.GetAjaxDistinctApi(&config))
}

func (h *DefaultApiHandler) GetAjaxDistinctApi(config *dbx.AjaxStructConfig) JsonAPIFunc {
	return func(c *gin.Context, j *jsonx.Json, us *session.UserSession) (ret interface{}, err error) {
		if !h.AuthAdminUserSession(us) {
			return nil, ErrAuthFailed
		}
		key, e := c.GetQuery("key")
		if !e {
			return
		}
		ret, err = config.GetStructFieldDistinct(key)
		return
	}
}

func (h *DefaultApiHandler) GetAjaxApi(config *dbx.AjaxStructConfig) JsonAPIFunc {
	return func(c *gin.Context, j *jsonx.Json, us *session.UserSession) (ret interface{}, err error) {
		if !h.AuthAdminUserSession(us) {
			return nil, ErrAuthFailed
		}
		ret, err = h.GetDataByAjaxQuery(c, j, us, config)
		return
	}
}
