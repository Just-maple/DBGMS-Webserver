package permission

import "reflect"

type Config struct {
	TableMap         map[string]Table
	PermissionConfig PermissionConfig
}

type TableConfig interface{}

type Table struct {
	FilesName    string
	TableData    []byte
	StructConfig *StructConfig
	Md5Hash      string
	TableConfig  TableConfig
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