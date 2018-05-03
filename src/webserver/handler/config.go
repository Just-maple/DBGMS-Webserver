package handler

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"syncx"
	. "webserver/args"
	pm "webserver/permission"
	. "webserver/utilsx"
	"gopkg.in/mgo.v2/bson"
	"webserver/dbx"
)

const (
	extensionJson = ".json"
)

type TableController struct {
	handler          *DefaultApiHandler
	PermissionConfig *pm.Config
	path             string
	collection       *dbx.Collection
}

func (h *DefaultApiHandler) InjectTableController(cfg *pm.PermissionConfig) (err error) {
	h.TableController = &TableController{
		handler: h,
		path:    h.GetTablePath(),
		PermissionConfig: &pm.Config{
			TableType: cfg.TableType,
			FieldType: cfg.FieldType,
		},
		collection: cfg.Collection,
	}
	if h.TableController.collection != nil {
		err = h.TableController.initAllConfigTableFromDatabaseCollection()
	} else {
		err = h.TableController.initAllConfigTableFromFiles()
	}
	h.TableController.registerAPI()
	return
}

func (c *TableController) registerAPI() {
	postRoute := c.handler.GetApiHandlersFromMethod(http.MethodPost)
	getRoute := c.handler.GetApiHandlersFromMethod(http.MethodGet)
	allPermissionApi := postRoute.MakeRegisterGroup(c.AuthAllPermission)
	allPermissionApi.RegisterDefaultAPI("saveAllConfig", c.SaveAllTableConfig)
	allPermissionApi.RegisterDefaultAPI("editTable", c.EditTable)
	getRoute.RegisterDefaultAPI("table", c.GetAllConfigTable, c.AuthAllPermission)
	postRoute.RegisterDefaultAPI("table", func(args *APIArgs) (ret interface{}, err error) {
		return c.GetConfigTableFromMString(args)
	})
}

func (c *TableController) AuthAllPermission(args *APIArgs) bool {
	access := c.handler.GetAccessConfig(args)
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

func (c *TableController) writeTableAndUpdateServerConfig(tableName, data string) (err error) {
	err = c.PermissionConfig.InitTableConfig([]byte(data), tableName)
	if err != nil {
		return
	}
	if c.collection != nil {
		_, err = c.collection.UpsertId(tableName, bson.M{"$set": bson.M{"data": data}})
	} else {
		file := c.path + tableName + extensionJson
		err = ioutil.WriteFile(file, []byte(data), 0600)
	}
	return
}

func (c *TableController) GetConfigTableFromMString(args *APIArgs) (tableConfig map[string]string, err error) {
	var storeBytes []byte
	var storeHash map[string]string
	storeHashString := args.JsonKey("m").MustString("{}")
	tableConfig = c.getConfigTableFromArgs(args)
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

func (c *TableController) GetAllConfigTable(args *APIArgs) (ret interface{}, err error) {
	return c.GetConfigTableMap(), nil
}

func (c *TableController) GetConfigTableMap() (ret map[string]string) {
	ret = make(map[string]string)
	pmConfig := c.PermissionConfig.TableMap
	for key := range pmConfig {
		ret[key] = string(pmConfig[key].TableData)
	}
	return
}

func (c *TableController) getConfigTableFromArgs(args *APIArgs) (ret map[string]string) {
	ret = make(map[string]string)
	access := c.handler.GetAccessConfig(args)
	pmConfig := c.PermissionConfig.TableMap
	mapLock := new(sync.Mutex)
	syncx.TraverseMapWithFunction(pmConfig, func(key string) {
		if access != nil && pmConfig[key].TableConfig.AuthTablePermission(access) {
			encodeKey := base64.StdEncoding.EncodeToString([]byte(key))
			encodeData := BytesToMd5String(pmConfig[key].TableData) + xdEncode(pmConfig[key].TableData)
			mapLock.Lock()
			ret[encodeKey] = encodeData
			mapLock.Unlock()
		}
	})
	return
}

func xdEncode(data []byte) string {
	str := base64.StdEncoding.EncodeToString([]byte(hex.EncodeToString(data)))
	str = str[0:4] + "NmU2MTZkNjUyMjN" + str[4:]
	return str
}

func (c *TableController) initTableConfigFromFileInfo(file *os.FileInfo) (err error) {
	if !strings.Contains((*file).Name(), extensionJson) {
		return
	}
	data, err := c.readTableConfigFromFile((*file).Name())
	if err != nil {
		return
	}
	tableName := strings.Replace((*file).Name(), extensionJson, "", 1)
	err = c.PermissionConfig.InitTableConfig(data, tableName)
	return
}

func (c *TableController) readTableConfigFromFile(fileName string) (data []byte, err error) {
	return ioutil.ReadFile(c.path + fileName)
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

type TableData []struct {
	Name string `bson:"_id"`
	Data string `bson:"data"`
}

func (c *TableController) initAllConfigTableFromDatabaseCollection() (err error) {
	var tableData TableData
	err = c.collection.Find(nil).All(&tableData)
	if err != nil {
		return
	}
	c.PermissionConfig.TableMap = make(map[string]*pm.Table, len(tableData))
	for i := range tableData {
		err = c.PermissionConfig.InitTableConfig([]byte(tableData[i].Data), tableData[i].Name)
		if err != nil {
			log.Fatal(err)
		}
	}
	return
}

func (c *TableController) initAllConfigTableFromFiles() (err error) {
	if !IsExist(c.path) {
		err = os.Mkdir(c.path, 0700)
		if err != nil {
			return
		}
	}
	tableFiles, err := ioutil.ReadDir(c.path)
	if err != nil {
		return
	}
	c.PermissionConfig.TableMap = make(map[string]*pm.Table, len(tableFiles))
	for i := range tableFiles {
		err = c.initTableConfigFromFileInfo(&(tableFiles[i]))
		if err != nil {
			log.Fatal(err)
		}
	}
	return
}
