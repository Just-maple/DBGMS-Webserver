package permission

import (
	"reflect"
)

type Config struct {
	TableMap  map[string]*Table
	TableType reflect.Type
	FieldType reflect.Type
}

type TableConfig interface {
	AuthTablePermission(AccessConfig) bool
}

type Table struct {
	Name         string
	TableData    []byte
	Md5Hash      string
	StructConfig *StructConfig
	TableConfig  TableConfig
	TableType    reflect.Type
	StructType   reflect.Type
}

type StructConfig map[string]FieldConfig

type StructFieldList []string

type FieldConfig interface {
	AuthFieldPermission(AccessConfig) bool
}

type AccessConfig interface {
	AuthAllPermission() bool
}

type PermissionConfig struct {
	TableType  reflect.Type
	FieldType  reflect.Type
}


func NewPermissionConfig(table TableConfig, field FieldConfig) *PermissionConfig {
	return &PermissionConfig{
		reflect.TypeOf(table).Elem(),
		reflect.TypeOf(field).Elem(),
	}
}
