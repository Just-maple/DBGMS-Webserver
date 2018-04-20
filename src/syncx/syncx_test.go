package syncx

import (
	"testing"
)

const testAppend = "fafssafasffsa"

func Test(t *testing.T) {
	var testMap = map[string]string{
		"test1": "",
		"test2": "",
	}
	err := TraverseMapWithFunction(testMap, func(key string) {
		testMap[key] = key + testAppend
	})
	if testMap["test1"] != "test1"+testAppend || testMap["test2"] != "test2"+testAppend || err != nil {
		t.Error("test failed", err)
	}
	
}

func BenchmarkTraverseSliceWithFunction(b *testing.B) {
	var testSlice []int
	for i := 0; i < 1000000; i++ {
		testSlice = append(testSlice, i)
	}
	err := TraverseSliceWithFunction(testSlice, func(i int) {
		testSlice[i] += i
	})
	if err != nil {
		b.Error(err)
	}
	
}
