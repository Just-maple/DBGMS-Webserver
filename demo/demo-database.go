package main

import (
	"webserver/dbx"
	"webserver/permission"
	"webserver/handler"
)

type DataBase struct {
	//database struct interface implement server.DB
	AnyCollection *dbx.Collection `collection:""`
	WXUser        *dbx.Collection
	//any public *mgo.Collection will init when database init
	//collection will init from Collection name in lower case like "anycollection" or tag collection
}

func (db *DataBase) GetAccessConfig(args *handler.APIArgs) permission.AccessConfig {
	//database struct implement auth super admin user
	//define your logic here
	return &SuperAdminAccess{args.Query("userid") == "User is Admin", args.Query("userid") == "User is Super"}
}

func (access *SuperAdminAccess) AuthTablePermission(config permission.TableConfig) bool {
	return (!config.PermissionConfig.GetTableConfig().(AdminConfig).NeedAdmin || access.isAdmin) && (!config.PermissionConfig.GetTableConfig().(AdminConfig).NeedSuperAdmin || access.isSuper)
}

type SuperAdminAccess struct {
	isAdmin bool
	isSuper bool
}

func (access *SuperAdminAccess) AuthAllPermission() bool {
	return access.isSuper
}
func (access *SuperAdminAccess) AuthFieldPermission(config permission.FieldConfig) bool {
	return (!config.(AdminStructConfig).SuperAdmin || access.isSuper) && (!config.(AdminStructConfig).Admin || access.isAdmin)
}
