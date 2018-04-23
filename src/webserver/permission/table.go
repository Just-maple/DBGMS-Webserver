package permission

import (
	"encoding/json"
	"reflect"
	"webserver/utilsx"
	"logger"
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

func (t *Table) InitTableConfig() (err error) {
	structTable, err := t.InitTableConfigMapFromBytes(t.TableData)
	if err != nil {
		return
	}
	var s = reflect.New(t.TableType).Interface()
	err = json.Unmarshal(t.TableData, s)
	if err != nil {
		return
	}
	t.TableConfig = reflect.ValueOf(s).Interface().(TableConfig)
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
				s := reflect.New(t.StructType).Interface()
				err = json.Unmarshal(tmp, s)
				if err != nil {
					return
				}
				mapLock.Lock()
				structConfig[key] = reflect.ValueOf(s).Interface().(FieldConfig)
				mapLock.Unlock()
			}
		})
	return
}

func (t *Config) InitTableConfig(data []byte, tableName string) (err error) {
	var config = &Table{
		FilesName:  tableName + extensionJson,
		TableData:  data,
		Md5Hash:    utilsx.BytesToMd5String(data),
		TableType:  t.TableType,
		StructType: t.FieldType,
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
