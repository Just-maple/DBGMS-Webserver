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

var (
	_ permission.AccessConfig = &SuperAdminAccess{}
)

func MakeSuperAdminAccess(IsAdmin, IsSuper bool) *SuperAdminAccess {
	return &SuperAdminAccess{
		IsAdmin: IsAdmin,
		IsSuper: IsSuper,
	}
}

func GetAdminPermissionConfig() *permission.PermissionConfig {
	return permission.NewPermissionConfig(new(AdminTableConfig), new(AdminStructConfig))
}

//define how your access config check by permission config
func (config *AdminTableConfig) AuthTablePermission(access permission.AccessConfig) bool {
	ac, valid := access.(*SuperAdminAccess)
	return valid && (!config.NeedAdmin || ac.IsAdmin) && (!config.NeedSuperAdmin || ac.IsSuper)
}
func (config *AdminStructConfig) AuthFieldPermission(access permission.AccessConfig) bool {
	ac, valid := access.(*SuperAdminAccess)
	return valid && (!config.SuperAdmin || ac.IsSuper) && (!config.Admin || ac.IsAdmin)
}

//define the all permission adjust
func (access *SuperAdminAccess) AuthAllPermission() bool {
	return access.IsSuper
}
