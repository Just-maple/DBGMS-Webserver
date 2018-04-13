package permission

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"reflect"
	"webserver/utilsx"
)

const (
	extensionJson        = ".json"
	privateKeyAdmin      = "_admin"
	privateKeySuperAdmin = "_superAdmin"
	privateKey           = "_"
)

type TableConfig struct {
	FilesName      string
	TableData      []byte
	StructConfig   *StructConfig
	Md5Hash        string
	NeedAdmin      bool
	NeedSuperAdmin bool
}

type TableConfigMap map[string]*TableConfig

func IsPrivateKey(key string) bool {
	return key[:1] == privateKey
}

func InitTableConfigMapFromBytes(data []byte) (res map[string]interface{}, err error) {
	err = json.Unmarshal(data, &res)
	return
}

func InitTableConfigFromBytes(data []byte) (fieldTableConfig StructConfig, needAdmin, needSuperAdmin bool, err error) {
	tablePermission, err := InitTableConfigMapFromBytes(data)
	if err != nil {
		return
	}
	for key := range tablePermission {
		if IsPrivateKey(key) {
			if key == privateKeyAdmin {
				needAdmin, _ = tablePermission[key].(bool)
			} else if key == privateKeySuperAdmin {
				needSuperAdmin, _ = tablePermission[key].(bool)
			}
			delete(tablePermission, key)
		}
	}
	tmp, err := json.Marshal(tablePermission)
	if err != nil {
		return
	}
	json.Unmarshal(tmp, &fieldTableConfig)
	if err != nil {
		return
	}
	return
}

func (t TableConfigMap) GetConfigTableFromContext(c *gin.Context) (config *StructConfig, has bool) {
	table, has := t[c.Param("api")]
	if has {
		config = table.StructConfig
	}
	return
}

func (t TableConfigMap) InitTableConfig(data []byte, tableName string) (err error) {
	fieldTableConfig, needAdmin, needSuperAdmin, err := InitTableConfigFromBytes(data)
	if err != nil {
		return
	}
	t[tableName] = &TableConfig{
		tableName + extensionJson,
		data,
		&fieldTableConfig,
		utilsx.BytesToMd5String(data),
		needAdmin,
		needSuperAdmin,
	}
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
