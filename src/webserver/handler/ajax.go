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

type ajaxResult struct {
	Data  interface{} `json:"data"`
	Count int         `json:"cnt"`
}

func (h *DefaultApiHandler) getDataByAjaxQuery(args *APIArgs, ajaxConfig *dbx.AjaxStructConfig) (res ajaxResult, err error) {
	query, err := h.getAjaxQuery(args)
	if err != nil {
		return
	}
	res.Data, res.Count, err = query.AjaxSearch(ajaxConfig)
	if err != nil {
		log.Error(err)
		return
	}
	access := h.GetAccessConfigFromArgs(args)
	res.Data = query.MakeAjaxReturnWithSelectKeysAndPermissionControl(res.Data, access)
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
