package dbx

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"webserver/utilsx"
)

type Collection struct {
	*mgo.Collection
}

func (c *Collection) Insert(docs interface{}) (err error) {
	return c.Collection.Insert(docs)
}

func (c *Collection) Pipe(pipelines interface{}) *mgo.Pipe {
	selector := []interface{}{
		bson.M{BsonSelectorMatch: BsonSelectUnDeleted},
	}
	v := reflect.ValueOf(pipelines)
	l := v.Len()
	for i := 0; i < l; i++ {
		selector = append(selector, v.Index(i).Interface())
	}
	return c.Collection.Pipe(selector)
}

func (c *Collection) Update(selector interface{}, update interface{}) (err error) {
	return c.Collection.Update(bson.M{
		BsonSelectorAnd: []interface{}{
			selector,
			BsonSelectUnDeleted,
		}}, update)
}

func (c *Collection) UpdateAll(selector interface{}, update interface{}) (info *mgo.ChangeInfo, err error) {
	return c.Collection.UpdateAll(bson.M{
		BsonSelectorAnd: []interface{}{
			selector,
			BsonSelectUnDeleted,
		}}, update)
}

func (c *Collection) UpdateId(id interface{}, update interface{}) (err error) {
	return c.Collection.Update(bson.M{
		BsonSelectorAnd: []interface{}{
			bson.M{FieldId: id},
			BsonSelectUnDeleted,
		}}, update)
}

func (c *Collection) Upsert(selector interface{}, update interface{}) (info *mgo.ChangeInfo, err error) {
	return c.Collection.Upsert(bson.M{
		BsonSelectorAnd: []interface{}{
			selector,
			BsonSelectUnDeleted,
		}}, update)
}

func (c *Collection) UpsertId(id interface{}, update interface{}) (info *mgo.ChangeInfo, err error) {
	return c.Collection.Upsert(bson.M{
		BsonSelectorAnd: []interface{}{
			bson.M{FieldId: id},
			BsonSelectUnDeleted,
		}}, update)
}

func (c *Collection) FindId(id interface{}) *mgo.Query {
	return c.Collection.Find(bson.M{
		BsonSelectorAnd: []interface{}{
			bson.M{FieldId: id},
			BsonSelectUnDeleted,
		}})
}

func (c *Collection) FindAll(query interface{}, in interface{}) error {
	return c.Find(query).All(in)
}

func (c *Collection) Find(query interface{}) *mgo.Query {
	if query == nil {
		query = bson.M{}
	}
	return c.Collection.Find(bson.M{
		BsonSelectorAnd: []interface{}{
			query,
			BsonSelectUnDeleted,
		}})
}

func (c *Collection) RemoveId(id interface{}) (err error) {
	err = c.Collection.UpdateId(id, BsonSetDeleted)
	return
}

func (c *Collection) Remove(selector interface{}) (err error) {
	err = c.Collection.Update(selector, BsonSetDeleted)
	return
}


func (c *Collection) GenerateRawStruct() (structRaw string, err error) {
	var d map[string]interface{}
	err = c.Find(nil).Sort("-t_create").One(&d)
	if err != nil {
		return
	}
	structRaw = utilsx.GenerateMapStruct(d, c.Name)
	return structRaw, err
}


