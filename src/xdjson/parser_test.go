package xdjson

import "testing"
import "encoding/json"

func TestInit(t *testing.T) {
	test := `{"sb":[11,{"tes":"asf" ,"asf":{"et":12412} },"242",{"tt":[1,[2,{"s":""}]],"etet":{"aa":[222,111],"fsf":{"wr":122,"22":2142}},"etetsas":[1324,15,"asf",false]}]}`
	testj := Init(test)
	var testmap = make(map[string]interface{})
	json.Unmarshal([]byte(test), &testmap)
	t.Log(testj.Map)
	t.Log(testmap)
}
