package main

import "webserver/dbx"

type DataBase struct {
	//database struct interface implement server.DB
	AnyCollection *dbx.Collection
	//any public *mgo.Collection will init when database init
	//collection will init from Collection name in lower case like "anycollection"
}

func (db *DataBase) AuthSuperAdminUser(userId string) (bool, bool) {
	//database struct implement auth super admin user
	//define your logic here
	return userId == "User is Admin", userId == "User is Super"
}

func (db *DataBase) AuthAdminUser(userId string) bool {
	//database struct implement auth  admin user
	//define your logic here
	return userId == "User is Admin"
}
