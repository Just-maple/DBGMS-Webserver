package permission

import (
	"reflect"
	"webserver/dbx"
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
	Collection *dbx.Collection
}

func (config *PermissionConfig) UseCollection(collection *dbx.Collection) {
	config.Collection = collection
}

func NewPermissionConfig(table TableConfig, field FieldConfig) *PermissionConfig {
	return &PermissionConfig{
		reflect.TypeOf(table).Elem(),
		reflect.TypeOf(field).Elem(),
		nil,
	}
}
