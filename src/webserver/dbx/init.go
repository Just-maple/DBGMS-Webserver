package dbx

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"strings"
	"time"
	"webserver/logger"
	"webserver/server"
)

var log = logger.Log
var BsonNotEqualTrue = bson.M{BsonSelectorNotEqual: true}
var BsonSelectUnDeleted = bson.M{FieldDeleted: BsonNotEqualTrue}
var BsonSetDeleted = bson.M{BsonSelectorSet: bson.M{FieldDeleted: true}}

func NewMgoDataBase(mgoURL, dbName string) (db *mgo.Database, err error) {
	dialInfo, err := mgo.ParseURL(mgoURL)
	if err != nil {
		log.Fatalf("Parse MgoDB Url err(%v)", err)
		return
	}
	
	dbSession, err := mgo.DialWithTimeout(mgoURL, time.Second*60)
	if err != nil {
		log.Errorf("Connect MgoDB Error = (%v)", err)
		return
	}
	dbSession.SetSafe(&mgo.Safe{})
	if dbName == "" {
		db = dbSession.DB(dialInfo.Database)
	} else {
		db = dbSession.DB(dbName)
	}
	return
}

func NewMgoDB(mgoURL string, db server.DB) (err error) {
	log.Debugf("Start Init MgoDB ,Url =  [ %v ]", mgoURL)
	if err != nil {
		return
	}
	t := reflect.TypeOf(db).Elem()
	s := reflect.ValueOf(db).Elem()
	cp := reflect.TypeOf(&Collection{})
	tFieldNum := t.NumField()
	newDB, err := NewMgoDataBase(mgoURL, "")
	if err != nil {
		return
	}
	for k := 0; k < tFieldNum; k++ {
		if t.Field(k).Type != cp {
			continue
		}
		fieldName := t.Field(k).Name
		if s.FieldByName(fieldName).CanSet() {
			var fieldCollection = &Collection{newDB.C(strings.ToLower(fieldName))}
			s.FieldByName(fieldName).Set(reflect.ValueOf(fieldCollection))
			//log.Debugf("Success Init Collection [ %v ]", fieldName)
		}
	}
	log.Debugf("Success Init MgoDB")
	return
}
