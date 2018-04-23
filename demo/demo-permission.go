package main

import (
	"webserver/args"
	"webserver/permission"
)

type AdminTableConfig struct {
	NeedAdmin      bool `json:"_admin"`
	NeedSuperAdmin bool `json:"_superAdmin"`
}

type AdminStructConfig struct {
	Admin      bool `json:"admin"`
	SuperAdmin bool `json:"superAdmin"`
}

type SuperAdminAccess struct {
	isAdmin bool
	isSuper bool
}

func (h *ApiHandler) GetPermissionConfig() *permission.PermissionConfig {
	return permission.NewPemissionConfig(new(AdminTableConfig), new(AdminStructConfig))
}

func (h *ApiHandler) GetAccessConfig(args *args.APIArgs) permission.AccessConfig {
	//database struct implement auth super admin user
	//define your logic here
	return &SuperAdminAccess{
		args.Query("userid") == "admin userId",
		args.Query("userid") == "super admin userId",
	}
}

func (config *AdminTableConfig) AuthTablePermission(access permission.AccessConfig) bool {
	return (!config.NeedAdmin || access.(*SuperAdminAccess).isAdmin) && (!config.NeedSuperAdmin || access.(*SuperAdminAccess).isSuper)
}
func (config *AdminStructConfig) AuthFieldPermission(access permission.AccessConfig) bool {
	return (!config.SuperAdmin || access.(*SuperAdminAccess).isSuper) && (!config.Admin || access.(*SuperAdminAccess).isAdmin)
}

func (access *SuperAdminAccess) AuthAllPermission() bool {
	return access.isSuper
}
