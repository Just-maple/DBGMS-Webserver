package access

import (
	"webserver/permission"
)

//define table permission config
type AdminTableConfig struct {
	NeedAdmin      bool `json:"_admin"`
	NeedSuperAdmin bool `json:"_superAdmin"`
}

//define struct permission config
type AdminStructConfig struct {
	Admin      bool `json:"admin"`
	SuperAdmin bool `json:"superAdmin"`
}

//define access config
//this access config will return with your access config function
//and you can check it with table or struct permission config
//then return a bool value to decide it can get by user or not
type SuperAdminAccess struct {
	IsAdmin bool
	IsSuper bool
}

//define how your access config check by permission config
func (config *AdminTableConfig) AuthTablePermission(access permission.AccessConfig) bool {
	return (!config.NeedAdmin || access.(*SuperAdminAccess).IsAdmin) && (!config.NeedSuperAdmin || access.(*SuperAdminAccess).IsSuper)
}
func (config *AdminStructConfig) AuthFieldPermission(access permission.AccessConfig) bool {
	return (!config.SuperAdmin || access.(*SuperAdminAccess).IsSuper) && (!config.Admin || access.(*SuperAdminAccess).IsAdmin)
}

//define the all permission adjust
func (access *SuperAdminAccess) AuthAllPermission() bool {
	return access.IsSuper
}
