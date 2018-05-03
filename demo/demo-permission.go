package main

import (
	"access"
	"webserver/permission"
	"webserver/args"
	"webserver/user"
	"gopkg.in/mgo.v2/bson"
)

func (h *ApiHandler) GetPermissionConfig() *permission.PermissionConfig {
	return permission.NewPemissionConfig(new(access.AdminTableConfig), new(access.AdminStructConfig))
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
	return &access.SuperAdminAccess{
		userdata.Level == 0,
		userdata.Level == 0,
	}
}