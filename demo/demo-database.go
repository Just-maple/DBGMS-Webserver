package main

import (
	"webserver/dbx"
	"webserver/permission"
)

type DataBase struct {
	//database struct interface implement server.DB
	AnyCollection *dbx.Collection `collection:""`
	//any public *mgo.Collection will init when database init
	//collection will init from Collection name in lower case like "anycollection" or tag collection
}

func (db *DataBase) GetAccessConfig(userId string) (permission.AccessConfig) {
	//database struct implement auth super admin user
	//define your logic here
	return &SuperAdminAccesss{userId == "User is Admin", userId == "User is Super"}
}

func (access *SuperAdminAccesss) AuthTablePermission(config *permission.TableConfig) bool {
	return (!config.NeedAdmin || access.isAdmin) && (!config.NeedSuperAdmin || access.isSuper)
}

type SuperAdminAccesss struct {
	isAdmin bool
	isSuper bool
}

func (access *SuperAdminAccesss) AuthAllPermission() bool {
	return access.isSuper
}
func (access *SuperAdminAccesss) AuthPermission(config *permission.StructFieldConfig) bool {
	return (!config.SuperAdmin || access.isSuper) && (!config.Admin || access.isAdmin)
}
