package table

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"syncx"
	. "webserver/args"
	. "webserver/utilsx"
)

const (
	extensionJson = ".json"
)

func (c *Controller) writeTableAndUpdateServerConfig(tableName, data string) (err error) {
	err = c.PermissionConfig.InitTableConfig([]byte(data), tableName)
	if err != nil {
		return
	}
	if c.Collection != nil {
		_, err = c.Collection.UpsertId(tableName, bson.M{"$set": bson.M{"data": data}})
	} else {
		file := c.Path + tableName + extensionJson
		err = ioutil.WriteFile(file, []byte(data), 0600)
	}
	return
}

func (c *Controller) GetConfigTableFromArgs(args *APIArgs) (tableConfig map[string]string, err error) {
	var storeBytes []byte
	var storeHash map[string]string
	storeHashString := args.JsonKey("m").MustString("{}")
	tableConfig = c.getConfigTableFromArgs(args)
	storeBytes, err = base64.StdEncoding.DecodeString(storeHashString)
	if err != nil {
		log.Error(err)
		return
	}
	json.Unmarshal(storeBytes, &storeHash)
	for key := range tableConfig {
		if storeHash[key] == tableConfig[key][:32] {
			tableConfig[key] = storeHash[key]
		}
	}
	return
}

func (c *Controller) GetConfigTableMap() (ret map[string]string) {
	ret = make(map[string]string)
	pmConfig := c.PermissionConfig.TableMap
	for key := range pmConfig {
		ret[key] = string(pmConfig[key].TableData)
	}
	return
}

func (c *Controller) getConfigTableFromArgs(args *APIArgs) (ret map[string]string) {
	ret = make(map[string]string)
	access := c.OauthAccessConfig(args)
	v, _ := args.UserId()
	if !v {
		return nil
	}
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

func (c *Controller) initTableConfigFromFileInfo(file *os.FileInfo) (err error) {
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

func (c *Controller) readTableConfigFromFile(fileName string) (data []byte, err error) {
	return ioutil.ReadFile(c.Path + fileName)
}
