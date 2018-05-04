package table

import (
	. "webserver/args"
)

func (c *Controller) GetTableFromHashStore(args *APIArgs) (ret interface{}, err error) {
	return c.GetConfigTableFromArgs(args)
}

func (c *Controller) AuthAllPermission(args *APIArgs) bool {
	v, _ := args.UserId()
	if !v {
		return false
	}
	access := c.OauthAccessConfig(args)
	return access != nil && access.AuthAllPermission()
}

func (c *Controller) SaveAllTableConfig(args *APIArgs) (ret interface{}, err error) {
	tableMap := args.Json.MustMap()
	for key := range tableMap {
		err = c.writeTableAndUpdateServerConfig(key, tableMap[key].(string))
		if err != nil {
			break
		}
	}
	return
}

func (c *Controller) EditTable(args *APIArgs) (ret interface{}, err error) {
	tableName := args.JsonKey("name").MustString()
	data := args.JsonKey("table").MustString()
	err = c.writeTableAndUpdateServerConfig(tableName, data)
	return
}

func (c *Controller) GetAllConfigTable(args *APIArgs) (ret interface{}, err error) {
	return c.GetConfigTableMap(), nil
}
