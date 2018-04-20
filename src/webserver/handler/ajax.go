package handler

import (
	"webserver/dbx"
	. "webserver/args"
)

func (h DefaultApiHandler) GetAjaxQuery(args *APIArgs) (res *dbx.AjaxQuery, err error) {
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
	h.ApiGetHandlers.RegisterDefaultAPI(distinctApiAddr, h.GetAjaxDistinctApi(&config))
}

func (h *DefaultApiHandler) GetAjaxDistinctApi(config *dbx.AjaxStructConfig) DefaultAPIFunc {
	return func(args *APIArgs) (ret interface{}, err error) {
		if config.AuthCheck != nil && !config.AuthCheck(args) {
			err = ErrAuthFailed
			return
		}
		key := args.Query("key")
		if key == "" {
			return
		}
		ret, err = config.GetStructFieldDistinct(key)
		return
	}
}

func (h *DefaultApiHandler) GetAjaxApi(config *dbx.AjaxStructConfig) DefaultAPIFunc {
	return func(args *APIArgs) (ret interface{}, err error) {
		if config.AuthCheck != nil && !config.AuthCheck(args) {
			err = ErrAuthFailed
			return
		}
		ret, err = h.GetDataByAjaxQuery(args, config)
		return
	}
}
