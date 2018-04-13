package handler

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"os"
	"strings"
	"webserver/jsonx"
	pm "webserver/permission"
	"webserver/session"
	. "webserver/utilsx"
)

const (
	extensionJson = ".json"
)

func (h *DefaultApiHandler) SaveAllTableConfig(c *gin.Context, j *jsonx.Json, us *session.UserSession) (ret interface{}, err error) {
	access := h.AuthUserSession(us)
	if !access.AuthAllPermission() {
		return nil, ErrAuthFailed
	}
	tableMap := j.MustMap()
	for key := range tableMap {
		err = h.WriteTableAndUpdateServerConfig(key, tableMap[key].(string))
		if err != nil {
			break
		}
	}
	return
}

func (h *DefaultApiHandler) WriteTableAndUpdateServerConfig(tableName, data string) (err error) {
	file := h.config.GetTablePath() + tableName + extensionJson
	err = ioutil.WriteFile(file, []byte(data), 0600)
	if err != nil {
		return
	}
	err = h.PermissionConfig.InitTableConfig([]byte(data), tableName)
	return
}

func (h *DefaultApiHandler) EditTable(c *gin.Context, j *jsonx.Json, us *session.UserSession) (ret interface{}, err error) {
	access := h.AuthUserSession(us)
	if !access.AuthAllPermission() {
		return nil, ErrAuthFailed
	}
	tableName := j.Get("name").MustString()
	data := j.Get("table").MustString()
	err = h.WriteTableAndUpdateServerConfig(tableName, data)
	return
}

func (h *DefaultApiHandler) GetConfigTableFromMString(m, UserId string) (tableConfig map[string]string, err error) {
	var storeBytes []byte
	var storeHash map[string]string
	tableConfig = h.ReadAllConfigTableFromServerTableConfig(true, UserId)
	storeBytes, err = base64.StdEncoding.DecodeString(m)
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

func (h *DefaultApiHandler) ReadAllConfigTableFromServerTableConfig(encode bool, userId string) (ret map[string]string) {
	ret = make(map[string]string)
	var access pm.AccessConfig
	if encode {
		access = h.db.GetAccessConfig(userId)
	}
	for key := range h.PermissionConfig {
		if encode {
			if access.AuthTablePermission(h.PermissionConfig[key]) {
				encodeKey := base64.StdEncoding.EncodeToString([]byte(key))
				encodeData := BytesToMd5String(h.PermissionConfig[key].TableData) + XdEncode(h.PermissionConfig[key].TableData)
				ret[encodeKey] = encodeData
			}
		} else {
			ret[key] = string(h.PermissionConfig[key].TableData)
		}
		
	}
	return
}

func XdEncode(data []byte) string {
	str := base64.StdEncoding.EncodeToString([]byte(hex.EncodeToString(data)))
	str = str[0:4] + "NmU2MTZkNjUyMjN" + str[4:]
	return str
}

func (h *DefaultApiHandler) InitTableConfigFromFileInfo(file *os.FileInfo) (err error) {
	if !strings.Contains((*file).Name(), extensionJson) {
		return
	}
	data, err := h.ReadTableConfigFromFile((*file).Name())
	if err != nil {
		return
	}
	tableName := strings.Replace((*file).Name(), extensionJson, "", 1)
	err = h.PermissionConfig.InitTableConfig(data, tableName)
	return
}

func (h *DefaultApiHandler) ReadTableConfigFromFile(fileName string) (data []byte, err error) {
	return ioutil.ReadFile(h.config.GetTablePath() + fileName)
}

func (h *DefaultApiHandler) InitAllConfigTableFromFiles() (err error) {
	tableFiles, err := ioutil.ReadDir(h.config.GetTablePath())
	if err != nil {
		return
	}
	h.PermissionConfig = make(map[string]*pm.TableConfig, len(tableFiles))
	for i := range tableFiles {
		err = h.InitTableConfigFromFileInfo(&(tableFiles[i]))
		if err != nil {
			log.Fatal(err)
		}
	}
	return
}

func (h *DefaultApiHandler) GetAllConfigTable(c *gin.Context, j *jsonx.Json, us *session.UserSession) (ret interface{}, err error) {
	access := h.AuthUserSession(us)
	if !access.AuthAllPermission() {
		return nil, ErrAuthFailed
	}
	ret = h.ReadAllConfigTableFromServerTableConfig(false, "")
	return
}
