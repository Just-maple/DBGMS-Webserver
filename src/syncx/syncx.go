package syncx

import (
	"sync"
	"reflect"
	"webserver/logger"
	"github.com/pkg/errors"
)

var log = logger.Log

func TraverseMapWithFunction(Map interface{}, function interface{}) (err error) {
	var mt = reflect.ValueOf(Map)
	var ft = reflect.ValueOf(function)
	if mt.Kind() != reflect.Map {
		return errors.New("invalid Map")
	}
	if mt.Type().Key() != ft.Type().In(0) {
		return errors.New("invalid function")
	}
	var wg = new(sync.WaitGroup)
	var keys = mt.MapKeys()
	for i := range keys {
		wg.Add(1)
		go func(index reflect.Value) {
			defer wg.Done()
			ft.Call([]reflect.Value{index})
		}(keys[i])
		
	}
	wg.Wait()
	return
}
func TraverseSliceWithFunction(slice interface{}, function interface{}) (err error) {
	var mt = reflect.ValueOf(slice)
	var ft = reflect.ValueOf(function)
	if mt.Kind() != reflect.Slice {
		return errors.New("invalid slice")
	}
	if ft.Type().In(0).Kind() != reflect.Int {
		return errors.New("invalid function")
	}
	var wg = new(sync.WaitGroup)
	for i := 0; i < mt.Len(); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ft.Call([]reflect.Value{reflect.ValueOf(i)})
		}(i)
		
	}
	wg.Wait()
	return
}

func Test() {
	var testMap = map[string]string{
		"test1": "",
		"test2": "",
	}
	err := TraverseMapWithFunction(testMap, func(key string) {
		testMap[key] = key + "fafssafasffsa"
	})
	if err != nil {
		panic(err)
	}
	log.Debug(testMap)
}
