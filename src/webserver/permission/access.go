package permission

type AccessConfig interface {
	AuthPermission(*StructFieldConfig) bool
	AuthTablePermission(*TableConfig) bool
	AuthAllPermission() bool
}
