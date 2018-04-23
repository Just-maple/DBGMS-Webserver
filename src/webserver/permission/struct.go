package permission

import (
	"reflect"
	"syncx"
)

func (structConfig *StructConfig) InitTablePermissionFieldList(ret interface{}, config AccessConfig) StructFieldList {
	retType := reflect.TypeOf(ret).Elem()
	for retType.Kind() == reflect.Ptr {
		retType = retType.Elem()
	}
	if retType.Kind() != reflect.Struct && retType.Kind() != reflect.Interface {
		return StructFieldList{}
	}
	return structConfig.GetFieldList(retType, config)
}

func (structConfig *StructConfig) InitTablePermission(ret interface{}, config AccessConfig) (res interface{}) {
	return structConfig.InitTablePermissionFieldList(ret, config).MakeFieldFilterReturnWithFieldList(ret)
}

func (structConfig *StructConfig) GetFieldList(retType reflect.Type, access AccessConfig) (fieldList StructFieldList) {
	allField := GetAllFieldNameFrom(retType)
	for _, fn := range allField {
		_, valid := retType.FieldByName(fn)
		if valid {
			if access.AuthAllPermission() {
				fieldList = append(fieldList, fn)
			} else if tmp, has := (*structConfig)[fn]; has && tmp.AuthFieldPermission(access) {
				fieldList = append(fieldList, fn)
			}
		}
	}
	return
}

func (fieldList StructFieldList) MergeList(in StructFieldList) (out StructFieldList) {
	var tmpMap = make(map[string]string, len(in))
	for _, key := range in {
		tmpMap[key] = key
	}
	for _, key2 := range fieldList {
		_, h := tmpMap[key2]
		if h {
			out = append(out, key2)
		}
	}
	return
}

func (fieldList StructFieldList) MakeFieldFilterReturnWithFieldList(in interface{}) interface{} {
	v := reflect.ValueOf(in)
	if v.Kind() != reflect.Slice {
		panic("type of arr not slice")
	}
	l := v.Len()
	var retSlice = make([]interface{}, l)
	syncx.TraverseSliceWithFunction(retSlice, func(si int) {
		itemV := reflect.ValueOf(v.Index(si).Interface())
		for itemV.Kind() == reflect.Ptr {
			itemV = itemV.Elem()
		}
		var retMap = make(map[string]interface{}, len(fieldList))
		for _, field := range fieldList {
			if itemV.FieldByName(field).IsValid() {
				s := itemV.FieldByName(field).Interface()
				retMap[field] = &s
			}
		}
		retSlice[si] = retMap
	})
	return retSlice
}
