package dbx

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"strings"
	"time"
	"webserver/permission"
)

const FieldTimeCreate = "t_create"

type MgoSearchCollection interface {
	Pipe(interface{}) *mgo.Pipe
	Find(interface{}) *mgo.Query
}

type AjaxMgoDBSearcher struct {
	Collection    MgoSearchCollection
	KeySelector   bson.M
	SelectorKeys  []string
	BsonMatcher   []bson.M
	FieldMatcher  map[string]interface{}
	SortKey       string
	SortBson      bson.M
	SortFieldName string
	Result        interface{}
}

type AjaxQuery struct {
	TimeStart        time.Time
	TimeEnd          time.Time
	SortKey          string
	Reverse          string
	LimitCount       int
	SkipCount        int
	MatcherMap       map[string]interface{}
	SelectKeys       permission.StructFieldList
	PermissionConfig *permission.StructConfig
}

const FieldTagBson = "bson"

type AjaxStructConfig struct {
	StructSlice interface{}
	Collection  MgoSearchCollection
	MiddleWare  func(interface{}) error
}

func (config *AjaxStructConfig) GetStructFieldDistinct(key string) (ret interface{}, err error) {
	field, e := reflect.TypeOf(config.StructSlice).FieldByName(key)
	if e {
		key = field.Tag.Get(FieldTagBson)
	}
	err = config.Collection.Find(bson.M{}).Distinct(key, &ret)
	return
}

func (query *AjaxQuery) MakeAjaxReturnWithSelectKeysAndPermissionControl(in interface{}, config permission.AccessConfig) (res interface{}) {
	fieldList := query.PermissionConfig.InitTablePermissionFieldList(in, config)
	keys := query.SelectKeys
	var tmp permission.StructFieldList
	if len(keys) != 0 {
		tmp = keys.MergeList(fieldList)
	} else {
		tmp = fieldList
	}
	res = tmp.MakeFieldFilterReturnWithFieldList(in)
	return
}

func (query *AjaxQuery) newAjaxMgoDBSearcher(collection MgoSearchCollection, result interface{}) (cnt int, err error) {
	resultv := reflect.ValueOf(result)
	if resultv.Kind() != reflect.Ptr || resultv.Elem().Kind() != reflect.Slice {
		panic("result argument must be a slice address")
	}
	return query.getMgoSearch(collection, result)
}

func initAjaxMgoDBSearcher(query *AjaxQuery, collection MgoSearchCollection, result interface{}) (ams *AjaxMgoDBSearcher) {
	ams = &AjaxMgoDBSearcher{
		Collection:    collection,
		FieldMatcher:  query.MatcherMap,
		SelectorKeys:  query.SelectKeys,
		SortFieldName: query.SortKey,
		Result:        result,
	}
	ams.getSortKey()
	ams.getSortBson(query.Reverse)
	ams.makeBsonMatcher(query.TimeStart, query.TimeEnd)
	return
}

func (ams *AjaxMgoDBSearcher) makeKeySelector() {
	ams.KeySelector = make(bson.M, len(ams.SelectorKeys))
	for i := range ams.SelectorKeys {
		field, e := reflect.TypeOf(ams.Result).Elem().Elem().FieldByName(ams.SelectorKeys[i])
		if e {
			ams.KeySelector[strings.Split(field.Tag.Get(FieldTagBson), ",")[0]] = 1
		} else {
			ams.KeySelector[ams.SelectorKeys[i]] = 1
		}
	}
	return
}

func (ams *AjaxMgoDBSearcher) makeBsonMatcher(st, et time.Time) {
	ams.BsonMatcher = []bson.M{{
		FieldTimeCreate: bson.M{
			"$lte": et,
			"$gt":  st,
		}}}
	for key, value := range ams.FieldMatcher {
		field, e := reflect.TypeOf(ams.Result).Elem().Elem().FieldByName(key)
		if e {
			var v bson.M
			if s, isStr := value.(string); isStr && s == "all" {
				v = bson.M{}
			} else {
				v = bson.M{strings.Split(field.Tag.Get(FieldTagBson), ",")[0]: value}
			}
			ams.BsonMatcher = append(ams.BsonMatcher, v)
		}
	}
	return
}

func (ams *AjaxMgoDBSearcher) getSortKey() {
	field, e := reflect.TypeOf(ams.Result).Elem().Elem().FieldByName(ams.SortFieldName)
	if e {
		ams.SortKey = strings.Split(field.Tag.Get(FieldTagBson), ",")[0]
	}
	if ams.SortKey == "" {
		ams.SortKey = FieldTimeCreate
	}
	return
}

func (ams *AjaxMgoDBSearcher) getTotalCount() (count int, err error) {
	return ams.Collection.Find(bson.M{
		"$and": ams.BsonMatcher,
	}).Count()
}

func (ams *AjaxMgoDBSearcher) getSortBson(reverse string) {
	var reverseI = 0
	if reverse == "" {
		reverseI = 1
	} else {
		reverseI = -1
	}
	ams.SortBson = bson.M{ams.SortKey: reverseI}
}
func (query *AjaxQuery) getMgoSearch(collection MgoSearchCollection, result interface{}) (cnt int, err error) {
	ams := initAjaxMgoDBSearcher(query, collection, result)
	if len(ams.SelectorKeys) != 0 {
		ams.makeKeySelector()
		err = ams.Collection.Pipe(
			[]bson.M{
				{"$match": bson.M{"$and": ams.BsonMatcher}},
				{"$sort": ams.SortBson},
				{"$project": ams.KeySelector},
			}).All(ams.Result)
	} else {
		err = ams.Collection.Pipe(
			[]bson.M{
				{"$match": bson.M{"$and": ams.BsonMatcher}},
				{"$sort": ams.SortBson},
				{"$skip": query.SkipCount},
				{"$limit": query.LimitCount},
			}).All(ams.Result)
	}
	if err != nil {
		return
	}
	return ams.getTotalCount()
}

func (query *AjaxQuery) AjaxSearch(structConfig *AjaxStructConfig) (res interface{}, count int, err error) {
	collection := structConfig.Collection
	tmp := reflect.New(reflect.SliceOf(reflect.TypeOf(structConfig.StructSlice))).Interface()
	count, err = query.newAjaxMgoDBSearcher(collection, tmp)
	if err != nil {
		return
	}
	if structConfig.MiddleWare != nil {
		err = structConfig.MiddleWare(tmp)
		if err != nil {
			return
		}
	}
	res = reflect.ValueOf(tmp).Elem().Interface()
	return
}
