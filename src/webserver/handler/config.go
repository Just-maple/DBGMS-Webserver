package handler

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	pm "webserver/permission"
	. "webserver/utilsx"
)

const (
	extensionJson = ".json"
)

type TableController struct {
	handler          TableHandler
	PermissionConfig pm.Config
	path             string
}

func InjectTableController(h TableHandler, PermissionConfig pm.PermissionConfig) (c *TableController, err error) {
	c = &TableController{
		handler: h,
		path:    h.GetTablePath(),
	}
	c.RegisterAPI()
	return c, c.InitAllConfigTableFromFiles(PermissionConfig)
}

func (c *TableController) RegisterAPI() {
	postRoute := c.handler.GetApiHandlersFromMethod(http.MethodPost)
	getRoute := c.handler.GetApiHandlersFromMethod(http.MethodGet)
	allPermissionApi := postRoute.MakeRegisterGroup(c.AuthAllPermission)
	allPermissionApi.RegisterDefaultAPI("saveAllConfig", c.SaveAllTableConfig)
	allPermissionApi.RegisterDefaultAPI("editTable", c.EditTable)
	getRoute.RegisterDefaultAPI("table", c.GetAllConfigTable, c.AuthAllPermission)
}

func (c *TableController) AuthAllPermission(args *APIArgs) bool {
	access := c.handler.GetAccessConfig(args)
	return access.AuthAllPermission()
}

func (c *TableController) SaveAllTableConfig(args *APIArgs) (ret interface{}, err error) {
	tableMap := args.Json.MustMap()
	for key := range tableMap {
		err = c.WriteTableAndUpdateServerConfig(key, tableMap[key].(string))
		if err != nil {
			break
		}
	}
	return
}

func (c *TableController) GetAllConfigTable(args *APIArgs) (ret interface{}, err error) {
	ret = c.ReadAllConfigTable()
	return
}

func (c *TableController) EditTable(args *APIArgs) (ret interface{}, err error) {
	tableName := args.JsonKey("name").MustString()
	data := args.JsonKey("table").MustString()
	err = c.WriteTableAndUpdateServerConfig(tableName, data)
	return
}

func (c *TableController) WriteTableAndUpdateServerConfig(tableName, data string) (err error) {
	file := c.path + tableName + extensionJson
	err = ioutil.WriteFile(file, []byte(data), 0600)
	if err != nil {
		return
	}
	err = c.PermissionConfig.InitTableConfig([]byte(data), tableName)
	return
}

func (c *TableController) GetConfigTableFromMString(args *APIArgs) (tableConfig map[string]string, err error) {
	var storeBytes []byte
	var storeHash map[string]string
	storeHashString := args.JsonKey("m").MustString("{}")
	tableConfig = c.ReadAllConfigTableFromServerTableConfig(args)
	storeBytes, err = base64.StdEncoding.DecodeString(storeHashString)
	if err != nil {
		log.Error(err)
		return
	}
	err = json.Unmarshal(storeBytes, &storeHash)
	if err != nil {
		log.Error(err)
		return
	}
	for key := range tableConfig {
		if storeHash[key] == tableConfig[key][:32] {
			tableConfig[key] = storeHash[key]
		}
	}
	return
}

func (c *TableController) ReadAllConfigTable() (ret map[string]string) {
	ret = make(map[string]string)
	pmConfig := c.PermissionConfig.TableMap
	for key := range pmConfig {
		ret[key] = string(pmConfig[key].TableData)
	}
	return
}

func (c *TableController) ReadAllConfigTableFromServerTableConfig(args *APIArgs) (ret map[string]string) {
	ret = make(map[string]string)
	access := c.handler.GetAccessConfig(args)
	pmConfig := c.PermissionConfig.TableMap
	for key := range pmConfig {
		if access.AuthTablePermission(pmConfig[key].TableConfig) {
			encodeKey := base64.StdEncoding.EncodeToString([]byte(key))
			encodeData := BytesToMd5String(pmConfig[key].TableData) + XdEncode(pmConfig[key].TableData)
			ret[encodeKey] = encodeData
		}
		
	}
	return
}

func XdEncode(data []byte) string {
	str := base64.StdEncoding.EncodeToString([]byte(hex.EncodeToString(data)))
	str = str[0:4] + "NmU2MTZkNjUyMjN" + str[4:]
	return str
}

func (c *TableController) InitTableConfigFromFileInfo(file *os.FileInfo) (err error) {
	if !strings.Contains((*file).Name(), extensionJson) {
		return
	}
	data, err := c.ReadTableConfigFromFile((*file).Name())
	if err != nil {
		return
	}
	tableName := strings.Replace((*file).Name(), extensionJson, "", 1)
	err = c.PermissionConfig.InitTableConfig(data, tableName)
	return
}

func (c *TableController) ReadTableConfigFromFile(fileName string) (data []byte, err error) {
	return ioutil.ReadFile(c.path + fileName)
}

func (c *TableController) InitAllConfigTableFromFiles(PermissionConfig pm.PermissionConfig) (err error) {
	tableFiles, err := ioutil.ReadDir(c.path)
	if err != nil {
		return
	}
	c.PermissionConfig.PermissionConfig = PermissionConfig
	c.PermissionConfig.TableMap = make(map[string]pm.Table, len(tableFiles))
	for i := range tableFiles {
		err = c.InitTableConfigFromFileInfo(&(tableFiles[i]))
		if err != nil {
			log.Fatal(err)
		}
	}
	return
}
