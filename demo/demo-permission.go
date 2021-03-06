package main

import (
	"access"
	"gopkg.in/mgo.v2/bson"
	"webserver/args"
	"webserver/controller/user"
	"webserver/permission"
)

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
	//define your access adjustment logic
	return access.MakeSuperAdminAccess(userdata.Level == 0, userdata.NickName == "admin", userid)
}
