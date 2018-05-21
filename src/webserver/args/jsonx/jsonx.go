package jsonx

import (
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"io"
	"reflect"
	"webserver/errorx"
)

const (
	JsonKeyId = "Id"
)

type Json struct {
	*simplejson.Json
}

func NewFromReader(r io.Reader) (*Json, error) {
	j, err := simplejson.NewFromReader(r)
	return &Json{j}, err
}
func New() *Json {
	return &Json{simplejson.New()}
}

func (j *Json) CallMethodByInstance(method interface{}, in interface{}) (err error) {
	if reflect.TypeOf(method).Kind() != reflect.Func || reflect.TypeOf(method).In(0) != reflect.TypeOf(in) {
		err = errorx.ErrMethodInvalid
		return
	}
	if err = j.Unmarshal(in); err != nil {
		return
	}
	res := reflect.ValueOf(method).Call([]reflect.Value{reflect.ValueOf(in)})
	reflect.ValueOf(&err).Elem().Set(res[0])
	return
}

func (j *Json) GetString(key string) string {
	return j.Get(key).MustString()
}

func (j *Json) Get(key string) *Json {
	return &Json{j.Json.Get(key)}
}

func (j *Json) GetStringId() string {
	return j.GetString(JsonKeyId)
}

func (j *Json) Unmarshal(in interface{}) (err error) {
	tmpB, err := json.Marshal(j)
	if err != nil {
		return
	}
	err = json.Unmarshal(tmpB, in)
	return
}
