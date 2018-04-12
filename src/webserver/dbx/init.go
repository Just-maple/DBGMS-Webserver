package dbx

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"strings"
	"time"
	"webserver/logger"
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

func NewMgoDB(mgoURL string, db interface{}) (err error) {
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
		collection := t.Field(k).Tag.Get("collection")
		field := s.FieldByName(fieldName)
		if field.CanSet() {
			if collection == "" {
				collection = strings.ToLower(fieldName)
			}
			var fieldCollection = &Collection{newDB.C(collection)}
			field.Set(reflect.ValueOf(fieldCollection))
			//log.Debugf("Success Init Collection [ %v ]", fieldName)MiddleWare
		}
	}
	log.Debugf("Success Init MgoDB")
	return
}
