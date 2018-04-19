package permission

import (
	"encoding/json"
	"reflect"
	"webserver/utilsx"
	"webserver/logger"
)

var log = logger.Log

const (
	extensionJson = ".json"
	privateKey    = "_"
)

type TableConfig struct {
	FilesName        string
	TableData        []byte
	StructConfig     *StructConfig
	Md5Hash          string
	PermissionConfig PermissionConfig
}

type TableMapConfig struct {
	TableMap         map[string]TableConfig
	PermissionConfig PermissionConfig
}

type PermissionConfig interface {
	GetTableConfig() (interface{})
	GetFieldConfig() (interface{})
}

func IsPrivateKey(key string) bool {
	return key[:1] == privateKey
}

func InitTableConfigMapFromBytes(data []byte) (res map[string]interface{}, err error) {
	err = json.Unmarshal(data, &res)
	return
}

func (t *TableConfig) InitTableConfig() (err error) {
	structTable, err := InitTableConfigMapFromBytes(t.TableData)
	if err != nil {
		return
	}
	var tmp = make(StructConfig, len(structTable))
	t.StructConfig = &tmp
	for key := range structTable {
		if IsPrivateKey(key) {
			delete(structTable, key)
		} else {
			tmp, err := json.Marshal(structTable[key])
			if err != nil {
				continue
			}
			s := reflect.New(reflect.TypeOf(t.PermissionConfig.GetFieldConfig())).Interface()
			err = json.Unmarshal(tmp, s)
			if err != nil {
				continue
			}
			(*t.StructConfig)[key] = reflect.ValueOf(s).Elem().Interface()
		}
	}
	if err != nil {
		panic(err)
	}
	return
}

func (t TableMapConfig) InitTableConfig(data []byte, tableName string) (err error) {
	var config = TableConfig{
		tableName + extensionJson,
		data,
		nil,
		utilsx.BytesToMd5String(data),
		t.PermissionConfig,
	}
	if err = config.InitTableConfig(); err != nil {
		return
	}
	t.TableMap[tableName] = config
	return
}

func GetAllFieldNameFrom(retType reflect.Type) (list []string) {
	for fi := 0; fi < retType.NumField(); fi++ {
		if retType.Field(fi).Anonymous {
			fit := retType.Field(fi).Type
			fitList := GetAllFieldNameFrom(fit)
			for fii := range fitList {
				list = append(list, fitList[fii])
			}
		} else {
			list = append(list, retType.Field(fi).Name)
		}
	}
	return
}
