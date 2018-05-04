package table

import (
	"net/http"
	. "webserver/args"
)

func (c *TableController) RegisterAPI() {
	allPermissionApi := c.DefaultController.MakeRegisterGroupByMethod(http.MethodPost, c.AuthAllPermission)
	getRoute := c.MakeRegisterGroupByMethod(http.MethodGet, c.AuthAllPermission)
	allPermissionApi.RegisterDefaultAPI("saveAllConfig", c.SaveAllTableConfig)
	allPermissionApi.RegisterDefaultAPI("editTable", c.EditTable)
	getRoute.RegisterDefaultAPI("table", c.GetAllConfigTable, )
	c.RegisterPostApi("table", c.GetTableFromHashStore)
}

func (c *TableController) GetTableFromHashStore(args *APIArgs) (ret interface{}, err error) {
	return c.GetConfigTableFromArgs(args)
}

func (c *TableController) AuthAllPermission(args *APIArgs) bool {
	access := c.GetAccessConfig(args)
	return access != nil && access.AuthAllPermission()
}

func (c *TableController) SaveAllTableConfig(args *APIArgs) (ret interface{}, err error) {
	tableMap := args.Json.MustMap()
	for key := range tableMap {
		err = c.writeTableAndUpdateServerConfig(key, tableMap[key].(string))
		if err != nil {
			break
		}
	}
	return
}

func (c *TableController) EditTable(args *APIArgs) (ret interface{}, err error) {
	tableName := args.JsonKey("name").MustString()
	data := args.JsonKey("table").MustString()
	err = c.writeTableAndUpdateServerConfig(tableName, data)
	return
}

func (c *TableController) GetAllConfigTable(args *APIArgs) (ret interface{}, err error) {
	return c.GetConfigTableMap(), nil
}
