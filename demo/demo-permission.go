package main

import (
	"webserver/args"
	"webserver/permission"
	"gopkg.in/mgo.v2/bson"
	"webserver/user"
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
	isAdmin bool
	isSuper bool
}

func (h *ApiHandler) GetPermissionConfig() *permission.PermissionConfig {
	return permission.NewPemissionConfig(new(AdminTableConfig), new(AdminStructConfig))
}

//define how user auth by your access config
//and make a access to check by permission config
func (h *ApiHandler) GetAccessConfig(args *args.APIArgs) permission.AccessConfig {
	//database struct implement auth super admin user
	//define your logic here
	_, userid := args.UserId()
	var userdata user.DefaultUser
	err := h.db.WXUser.FindId(bson.ObjectIdHex(userid)).One(&userdata)
	if err != nil {
		return nil
	}
	return &SuperAdminAccess{
		userdata.Level == 0,
		userdata.Level == 0,
	}
}

//define how your access config check by permission config
func (config *AdminTableConfig) AuthTablePermission(access permission.AccessConfig) bool {
	return (!config.NeedAdmin || access.(*SuperAdminAccess).isAdmin) && (!config.NeedSuperAdmin || access.(*SuperAdminAccess).isSuper)
}
func (config *AdminStructConfig) AuthFieldPermission(access permission.AccessConfig) bool {
	return (!config.SuperAdmin || access.(*SuperAdminAccess).isSuper) && (!config.Admin || access.(*SuperAdminAccess).isAdmin)
}

//define the all permission adjust
func (access *SuperAdminAccess) AuthAllPermission() bool {
	return access.isSuper
}
