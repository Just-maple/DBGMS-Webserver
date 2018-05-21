package main

import (
	"webserver/dbx"
)

type DataBase struct {
	dbx.Db
	//database struct interface implement server.DB
	AnyCollection   *dbx.Collection `collection:""`
	WXUser          *dbx.Collection
	PermissionTable *dbx.Collection `collection:"pm_table"`
	//any public *mgo.Collection will init when database init
	//collection will init from Collection name in lower case like "anycollection" or tag collection
}
