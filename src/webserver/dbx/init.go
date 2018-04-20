package dbx

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"strings"
	"time"
	"logger"
	"sync"
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
		log.Fatalf("Connect MgoDB Error = (%v)", err)
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
	var wg = new(sync.WaitGroup)
	for k := 0; k < tFieldNum; k++ {
		wg.Add(1)
		go func(k int) {
			defer wg.Done()
			if t.Field(k).Type != cp {
				return
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
			}
		}(k)
	}
	wg.Wait()
	log.Debugf("Success Init MgoDB")
	return
}
