package main

import (
	"webserver/dbx"
	"webserver/permission"
	"webserver/args"
)

type DataBase struct {
	//database struct interface implement server.DB
	AnyCollection *dbx.Collection `collection:""`
	WXUser        *dbx.Collection
	//any public *mgo.Collection will init when database init
	//collection will init from Collection name in lower case like "anycollection" or tag collection
}

func (db *DataBase) GetAccessConfig(args *args.APIArgs) permission.AccessConfig {
	//database struct implement auth super admin user
	//define your logic here
	return &SuperAdminAccess{args.Query("userid") == "User is Admin", args.Query("userid") == "User is Super"}
}

func (config *AdminConfig) AuthTablePermission(access permission.AccessConfig) bool {
	return (!config.NeedAdmin || access.(*SuperAdminAccess).isAdmin) && (!config.NeedSuperAdmin || access.(*SuperAdminAccess).isSuper)
}
func (config *AdminStructConfig) AuthFieldPermission(access permission.AccessConfig) bool {
	return (!config.SuperAdmin || access.(*SuperAdminAccess).isSuper) && (!config.Admin || access.(*SuperAdminAccess).isAdmin)
}

type SuperAdminAccess struct {
	isAdmin bool
	isSuper bool
}

func (access *SuperAdminAccess) AuthAllPermission() bool {
	return access.isSuper
}
