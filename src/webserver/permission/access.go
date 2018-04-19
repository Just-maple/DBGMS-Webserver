package permission

type AccessConfig interface {
	AuthFieldPermission(FieldConfig) bool
	AuthTablePermission(TableConfig) bool
	AuthAllPermission() bool
}
