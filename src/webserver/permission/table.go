package permission

import (
	"encoding/json"
	"logger"
	"reflect"
	"webserver/utilsx"
)

var log = logger.Log

const (
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
	for key := range structTable {
		if !IsPrivateKey(key) {
			tmp, err := json.Marshal(structTable[key])
			if err != nil {
				break
			}
			s := reflect.New(t.StructType).Interface()
			err = json.Unmarshal(tmp, s)
			if err != nil {
				break
			}
			structConfig[key] = reflect.ValueOf(s).Interface().(FieldConfig)
		}
	}
	t.StructConfig = &structConfig
	return
}

func (t *Config) InitTableConfig(data []byte, tableName string) (err error) {
	var config = &Table{
		Name:  tableName,
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
