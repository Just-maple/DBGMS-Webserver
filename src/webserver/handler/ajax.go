package handler

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
	"webserver/dbx"
	"webserver/jsonx"
	"webserver/session"
)

func TransAjaxQuery(args *APIArgs) (matcherMap map[string]interface{}, keys []string, skipCnt, limitCnt int, sortKey, reverse string, tSTime, tETime time.Time, err error) {
	keys = args.JsonKey("keys").MustStringArray()
	matcherMap = args.JsonKey("matcher").MustMap()
	
	skip := args.context.DefaultQuery("skip", "0")
	skipCnt, err = strconv.Atoi(skip)
	if err != nil {
		skipCnt = 0
		err = nil
	}
	limit := args.context.DefaultQuery("limit", "10")
	limitCnt, err = strconv.Atoi(limit)
	if err != nil {
		limitCnt = 0
		err = nil
	}
	sortKey = args.context.DefaultQuery("sort", "")
	reverse = args.context.DefaultQuery("reverse", "")
	tSTime, tETime = args.Time()
	return
}

func (h DefaultApiHandler) GetAjaxQuery(args *APIArgs) (res *dbx.AjaxQuery, err error) {
	matcherMap, keys, skipCnt, limitCnt, sortKey, reverse, tSTime, tETime, err := TransAjaxQuery(args)
	pmConfig, _ := args.GetConfigTable(&h.TableController.PermissionConfig)
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

func (h *DefaultApiHandler) GetDataByAjaxQuery(args *APIArgs, ajaxConfig *dbx.AjaxStructConfig) (res map[string]interface{}, err error) {
	query, err := h.GetAjaxQuery(args)
	if err != nil {
		return
	}
	access := h.GetAccessConfigFromArgs(args)
	data, count, err := query.AjaxSearch(ajaxConfig)
	if err != nil {
		log.Error(err)
		return
	}
	data = query.MakeAjaxReturnWithSelectKeysAndPermissionControl(data, access)
	res = map[string]interface{}{
		"data": data,
		"cnt":  count,
	}
	return
}

func (h *DefaultApiHandler) RegisterAjaxJsonApi(dataApiAddr, distinctApiAddr string, configMaker func() dbx.AjaxStructConfig) {
	config := configMaker()
	h.ApiPostHandlers.RegisterDefaultAPI(dataApiAddr, h.GetAjaxApi(&config))
	h.ApiGetHandlers.RegisterAPI(distinctApiAddr, h.GetAjaxDistinctApi(&config))
}

func (h *DefaultApiHandler) GetAjaxDistinctApi(config *dbx.AjaxStructConfig) JsonAPIFunc {
	return func(c *gin.Context, j *jsonx.Json, us *session.UserSession) (ret interface{}, err error) {
		valid, userId := us.AuthUserSession()
		if !valid || (config.AuthCheck != nil && !config.AuthCheck(userId)) {
			err = ErrAuthFailed
			return
		}
		key, e := c.GetQuery("key")
		if !e {
			return
		}
		ret, err = config.GetStructFieldDistinct(key)
		return
	}
}

func (h *DefaultApiHandler) GetAjaxApi(config *dbx.AjaxStructConfig) DefaultAPIFunc {
	return func(args *APIArgs) (ret interface{}, err error) {
		valid, userId := args.UserId()
		if !valid || (config.AuthCheck != nil && !config.AuthCheck(userId)) {
			err = ErrAuthFailed
			return
		}
		ret, err = h.GetDataByAjaxQuery(args, config)
		return
	}
}
