package permission

import "reflect"

type Config struct {
	TableMap  map[string]*Table
	TableType reflect.Type
	FieldType reflect.Type
}

type TableConfig interface{}

type Table struct {
	FilesName    string
	TableData    []byte
	Md5Hash      string
	StructConfig *StructConfig
	TableConfig  TableConfig
	TableType    reflect.Type
	StructType   reflect.Type
}

type StructConfig map[string]FieldConfig

type StructFieldList []string

type FieldConfig interface{}

type AccessConfig interface {
	AuthFieldPermission(FieldConfig) bool
	AuthTablePermission(TableConfig) bool
	AuthAllPermission() bool
}

type PermissionConfig interface {
	GetTableConfig() (reflect.Type)
	GetFieldConfig() (reflect.Type)
}
