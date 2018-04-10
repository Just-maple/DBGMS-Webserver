package permission

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"reflect"
	"sync"
)

const (
	extensionJson        = ".json"
	privateKeyAdmin      = "_admin"
	privateKeySuperAdmin = "_superAdmin"
	privateKey           = "_"
)

type StructFieldConfig struct {
	Admin      bool   `json:"admin"`
	SuperAdmin bool   `json:"superAdmin"`
	Name       string `json:"name"`
}

type ApiPermissionConfig struct {
	Admin      bool   `json:"admin"`
	SuperAdmin bool   `json:"superAdmin"`
}

type ApiConfigMap map[string]ApiPermissionConfig


type TableConfig struct {
	FilesName      string
	TableData      []byte
	StructConfig   *StructConfig
	Md5Hash        string
	NeedAdmin      bool
	NeedSuperAmind bool
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
		BytesToMd5String(data),
		needAdmin,
		needSuperAdmin,
	}
	return
}

func BytesToMd5String(data []byte) string {
	md5Ctx := md5.New()
	md5Ctx.Write(data)
	return hex.EncodeToString(md5Ctx.Sum(nil))
}

func (config *TableConfig) AuthPermission(isAdmin, isSuper bool) bool {
	return (!config.NeedAdmin || isAdmin) && (!config.NeedSuperAmind || isSuper)
}

type StructConfig map[string]StructFieldConfig
type StructFieldList []string

func (config *StructFieldConfig) AuthPermission(isAdmin, isSuper bool) bool {
	return (!config.SuperAdmin || isSuper) && (!config.Admin || isAdmin)
}

func GetAllFieldNameFromStruct(retType reflect.Type) (list []string) {
	for fi := 0; fi < retType.NumField(); fi++ {
		if retType.Field(fi).Anonymous {
			fit := retType.Field(fi).Type
			fitList := GetAllFieldNameFromStruct(fit)
			for fii := range fitList {
				list = append(list, fitList[fii])
			}
		} else {
			list = append(list, retType.Field(fi).Name)
		}
	}
	return
}

func (structConfig *StructConfig) InitTablePermissionFieldList(ret interface{}, isAdmin, isSuper bool) StructFieldList {
	retType := reflect.TypeOf(ret).Elem()
	for retType.Kind() == reflect.Ptr {
		retType = retType.Elem()
	}
	if retType.Kind() != reflect.Struct && retType.Kind() != reflect.Interface {
		return StructFieldList{}
	}
	return structConfig.GetFieldList(retType, isAdmin, isSuper)
}

func (structConfig *StructConfig) InitTablePermission(ret interface{}, isAdmin, isSuper bool) (res interface{}) {
	return structConfig.InitTablePermissionFieldList(ret, isAdmin, isSuper).MakeFieldFilterReturnWithFieldList(ret)
}

func (structConfig *StructConfig) GetFieldList(retType reflect.Type, isAdmin, isSuper bool) (fieldList StructFieldList) {
	allField := GetAllFieldNameFromStruct(retType)
	for _, fn := range allField {
		_, valid := retType.FieldByName(fn)
		if valid {
			if isSuper {
				fieldList = append(fieldList, fn)
			} else if tmp, has := (*structConfig)[fn]; has && tmp.AuthPermission(isAdmin, isSuper) {
				fieldList = append(fieldList, fn)
			}
		}
	}
	return
}

func (fieldList StructFieldList) MergeList(in StructFieldList) (out StructFieldList) {
	var tmpMap = make(map[string]string, len(in))
	for _, key := range in {
		tmpMap[key] = key
	}
	for _, key2 := range fieldList {
		_, h := tmpMap[key2]
		if h {
			out = append(out, key2)
		}
	}
	return
}

func (fieldList StructFieldList) MakeFieldFilterReturnWithFieldList(in interface{}) interface{} {
	v := reflect.ValueOf(in)
	if v.Kind() != reflect.Slice {
		panic("type of arr not slice")
	}
	l := v.Len()
	var retSlice = make([]interface{}, l)
	var wg sync.WaitGroup
	for si := 0; si < l; si++ {
		wg.Add(1)
		go func(si int) {
			defer wg.Done()
			itemV := reflect.ValueOf(v.Index(si).Interface())
			for itemV.Kind() == reflect.Ptr {
				itemV = itemV.Elem()
			}
			var retMap = make(map[string]interface{}, len(fieldList))
			for _, field := range fieldList {
				if itemV.FieldByName(field).IsValid() {
					s := itemV.FieldByName(field).Interface()
					retMap[field] = &s
				}
			}
			retSlice[si] = retMap
		}(si)
	}
	wg.Wait()
	return retSlice
}
