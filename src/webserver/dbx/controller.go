package dbx

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"webserver/errorx"
	"webserver/args/jsonx"
)

type CollectionController struct {
	collection   *Collection
	model        reflect.Type
	idKeyField   string
	json2BsonMap map[string]string
}

const FieldTCreate = "t_create"

func (c *Collection) CreateController(in interface{}) *CollectionController {
	var newIn = reflect.TypeOf(in)
	for newIn.Kind() == reflect.Ptr {
		newIn = newIn.Elem()
	}
	if newIn.Kind() != reflect.Struct {
		panic("controller model must be struct")
	}
	fi := newIn.NumField()
	var MatchIdKeyField bool
	var idKeyField string
	var json2BsonMap = make(map[string]string, fi)
	for i := 0; i < fi; i++ {
		if newIn.Field(i).Tag.Get("bson") == "_id" {
			if !MatchIdKeyField {
				MatchIdKeyField = true
				idKeyField = newIn.Field(i).Name
			} else {
				panic(fmt.Sprintf("IdKey repeat error ,match %v and %v", newIn.Field(i).Name, idKeyField))
			}
		}
		var FieldKey = newIn.Field(i).Name
		var JsonKey = newIn.Field(i).Tag.Get("json")
		if JsonKey != "" {
			FieldKey = JsonKey
		}
		json2BsonMap[FieldKey] = newIn.Field(i).Tag.Get("bson")
	}
	if !MatchIdKeyField {
		panic("IdKey match error,ensure model has field with tag bson")
	}
	return &CollectionController{
		collection:   c,
		model:        newIn,
		idKeyField:   idKeyField,
		json2BsonMap: json2BsonMap,
	}
}

func (c *CollectionController) NewModelSlice() (ret interface{}) {
	return reflect.New(reflect.SliceOf(c.model)).Interface()
}

func (c *CollectionController) GetIdString(json *jsonx.Json) (Id string) {
	return json.Get(c.idKeyField).MustString()
}
func (c *CollectionController) NewModel() (ret interface{}) {
	return reflect.New(c.model).Interface()
}

func (c *CollectionController) UpdateByJson(json *jsonx.Json) (err error) {
	var jmap = json.MustMap()
	var updator = c.NewModel()
	var Id = c.GetIdString(json)
	if !bson.IsObjectIdHex(Id) {
		return errorx.ErrIdInvalid
	}
	err = json.Unmarshal(updator)
	uv := reflect.ValueOf(updator)
	if err != nil {
		return
	}
	var bmap = make(bson.M, len(jmap))
	for key := range jmap {
		bmap[c.json2BsonMap[key]] = uv.FieldByName(key).Interface()
	}
	err = c.collection.UpdateId(bson.ObjectIdHex(Id), bson.M{"$set": bmap})
	return
}

func (c *CollectionController) RemoveByJson(json *jsonx.Json) (err error) {
	Id := c.GetIdString(json)
	if !bson.IsObjectIdHex(Id) {
		return errorx.ErrIdInvalid
	}
	return c.RemoveById(bson.ObjectIdHex(Id))
}

func (c *CollectionController) RemoveById(id bson.ObjectId) (err error) {
	err = c.collection.RemoveId(id)
	return
}

func (c *CollectionController) GetById(id bson.ObjectId) (ret interface{}, err error) {
	var res = c.NewModel()
	err = c.collection.FindId(id).One(res)
	ret = res
	return
}

func (c *CollectionController) GetAll(query interface{}) (ret interface{}, err error) {
	var res = c.NewModelSlice()
	err = c.collection.FindAll(query, res)
	if err != nil {
		return
	}
	ret = res
	return
}

func (c *CollectionController) NewFromJson(json *jsonx.Json) (err error) {
	var newModel = reflect.New(c.model).Elem()
	err = json.Unmarshal(newModel.Interface())
	if err != nil {
		return
	}
	newModel.FieldByName(c.idKeyField).Set(reflect.ValueOf(bson.NewObjectId()))
	newBson, err := StructToBson(newModel.Interface())
	if err != nil {
		return
	}
	newBson[FieldTCreate] = bson.NewObjectId()
	err = c.collection.Insert(newBson)
	return
}

func StructToBson(in interface{}) (newBson bson.M, err error) {
	tmpBytes, err := bson.Marshal(in)
	if err != nil {
		return
	}
	err = bson.Unmarshal(tmpBytes, newBson)
	return
}
