package main

import (
	"webserver/dbx"
	"webserver/permission"
)

type DataBase struct {
	//database struct interface implement server.DB
	AnyCollection *dbx.Collection `collection:""`
	WXUser        *dbx.Collection
	//any public *mgo.Collection will init when database init
	//collection will init from Collection name in lower case like "anycollection" or tag collection
}

func (db *DataBase) GetAccessConfig(userId string) (permission.AccessConfig) {
	//database struct implement auth super admin user
	//define your logic here
	return &SuperAdminAccess{userId == "User is Admin", userId == "User is Super"}
}

func (access *SuperAdminAccess) AuthTablePermission(config *permission.TableConfig) bool {
	return (!config.NeedAdmin || access.isAdmin) && (!config.NeedSuperAdmin || access.isSuper)
}

type SuperAdminAccess struct {
	isAdmin bool
	isSuper bool
}

func (access *SuperAdminAccess) AuthAllPermission() bool {
	return access.isSuper
}
func (access *SuperAdminAccess) AuthPermission(config *permission.StructFieldConfig) bool {
	return (!config.SuperAdmin || access.isSuper) && (!config.Admin || access.isAdmin)
}
