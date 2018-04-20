package permission

import (
	"encoding/json"
	"reflect"
	"webserver/utilsx"
	"webserver/logger"
	"syncx"
	"sync"
)

var log = logger.Log

const (
	extensionJson = ".json"
	privateKey    = "_"
)

func IsPrivateKey(key string) bool {
	return key[:1] == privateKey
}

func (t *Table) InitTableConfigMapFromBytes(data []byte) (res map[string]interface{}, err error) {
	err = json.Unmarshal(data, &res)
	return
}

func (t *Table) InitTableConfig(PermissionConfig PermissionConfig) (err error) {
	structTable, err := t.InitTableConfigMapFromBytes(t.TableData)
	if err != nil {
		return
	}
	var s = reflect.New(PermissionConfig.GetTableConfig()).Interface()
	err = json.Unmarshal(t.TableData, s)
	if err != nil {
		return
	}
	t.TableConfig = reflect.ValueOf(s).Elem().Interface()
	var structConfig = make(StructConfig, len(structTable))
	t.StructConfig = &structConfig
	var mapLock = new(sync.RWMutex)
	err = syncx.TraverseMapWithFunction(
		structTable, func(key string) {
			if !IsPrivateKey(key) {
				tmp, err := json.Marshal(structTable[key])
				if err != nil {
					return
				}
				s := reflect.New(PermissionConfig.GetFieldConfig()).Interface()
				err = json.Unmarshal(tmp, s)
				if err != nil {
					return
				}
				mapLock.Lock()
				structConfig[key] = reflect.ValueOf(s).Elem().Interface()
				mapLock.Unlock()
			}
		})
	return
}

func (t *Config) InitTableConfig(data []byte, tableName string) (err error) {
	var config = Table{
		tableName + extensionJson,
		data,
		nil,
		utilsx.BytesToMd5String(data),
		nil,
	}
	if err = config.InitTableConfig(t.PermissionConfig); err != nil {
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
