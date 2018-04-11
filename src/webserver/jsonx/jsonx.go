package jsonx

import (
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"reflect"
	"webserver/errorx"
)

type Json struct {
	*simplejson.Json
}

func (j *Json) CallMethodByInstance(method interface{}, intf interface{}) (err error) {
	if reflect.TypeOf(method).Kind() != reflect.Func || reflect.TypeOf(method).In(0) != reflect.TypeOf(intf) {
		err = errorx.ErrMethodInvalid
		return
	}
	if err = j.Unmarshal(intf); err != nil {
		return
	}
	res := reflect.ValueOf(method).Call([]reflect.Value{reflect.ValueOf(intf)})
	reflect.ValueOf(&err).Elem().Set(res[0])
	return
}

func (j *Json) GetString(key string) string {
	return j.Get(key).MustString()
}

func (j *Json) GetStringId() string {
	return j.GetString(JsonKeyId)
}

func (j *Json) Unmarshal(intf interface{}) (err error) {
	tmpB, err := json.Marshal(j)
	if err != nil {
		return
	}
	err = json.Unmarshal(tmpB, intf)
	return
}
