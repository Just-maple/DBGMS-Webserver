package dbx

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"logger"
	"reflect"
	"strings"
	"sync"
	"time"
)

var log = logger.Log
var BsonNotEqualTrue = bson.M{BsonSelectorNotEqual: true}
var BsonSelectUnDeleted = bson.M{FieldDeleted: BsonNotEqualTrue}
var BsonSetDeleted = bson.M{BsonSelectorSet: bson.M{FieldDeleted: true}}

type Db struct {
	db        *mgo.Database
	structMap map[string]reflect.Type
}

func (db *Db) GetCollection(name string) *Collection {
	return &Collection{db.db.C(name)}
}

func (db *Db) SetDB(DB *mgo.Database) {
	db.db = DB
	return
}

func (db *Db) GetNewStruct(name string) (interface{}, bool) {
	st, h := db.structMap[name]
	if !h {
		return nil, h
	}
	return reflect.New(st).Interface(), h
}

func (db *Db) GetNewStructSlice(name string) (interface{}, bool) {
	st, h := db.structMap[name]
	if !h {
		return nil, h
	}
	return reflect.New(reflect.SliceOf(st)).Interface(), h
}
func (db *Db) RegisterStruct(name string, in interface{}) () {
	if db.structMap == nil {
		db.structMap = make(map[string]reflect.Type)
	}
	_, h := db.structMap[name]
	if h {
		panic("key already use")
	}
	db.structMap[name] = reflect.TypeOf(in)
	return
}

type DB interface {
	GetCollection(string) *Collection
	SetDB(DB *mgo.Database)
	GetNewStruct(name string) (interface{}, bool)
	GetNewStructSlice(name string) (interface{}, bool)
	RegisterStruct(name string, in interface{}) ()
}

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

func NewMgoDB(mgoURL string, db DB) (err error) {
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
	db.SetDB(newDB)
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
