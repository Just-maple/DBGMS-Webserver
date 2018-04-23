package handler

import (
	. "webserver/args"
	"webserver/dbx"
)

func (h DefaultApiHandler) getAjaxQuery(args *APIArgs) (res *dbx.AjaxQuery, err error) {
	matcherMap, keys, skipCnt, limitCnt, sortKey, reverse, tSTime, tETime, err := args.TransAjaxQuery()
	pmConfig, _ := args.GetConfigTable(h.TableController.PermissionConfig)
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

func (h *DefaultApiHandler) getDataByAjaxQuery(args *APIArgs, ajaxConfig *dbx.AjaxStructConfig) (res map[string]interface{}, err error) {
	query, err := h.getAjaxQuery(args)
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

func (h *DefaultApiHandler) RegisterAjaxJsonApi(dataApiAddr, distinctApiAddr string, config *dbx.AjaxStructConfig, pm ...PermissionAuth) {
	h.ApiPostHandlers.RegisterDefaultAPI(dataApiAddr, h.getAjaxApi(config), pm...)
	h.ApiGetHandlers.RegisterDefaultAPI(distinctApiAddr, h.getAjaxDistinctApi(config), pm...)
}

func (h *DefaultApiHandler) getAjaxDistinctApi(config *dbx.AjaxStructConfig) DefaultAPIFunc {
	return func(args *APIArgs) (ret interface{}, err error) {
		key := args.Query("key")
		if key == "" {
			return
		}
		ret, err = config.GetStructFieldDistinct(key)
		return
	}
}

func (h *DefaultApiHandler) getAjaxApi(config *dbx.AjaxStructConfig) DefaultAPIFunc {
	return func(args *APIArgs) (ret interface{}, err error) {
		ret, err = h.getDataByAjaxQuery(args, config)
		return
	}
}
