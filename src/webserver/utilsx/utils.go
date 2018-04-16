package utilsx

import (
	"archive/zip"
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"os"
	"strconv"
	"time"
	"strings"
	"reflect"
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

func TransTime(c *gin.Context) (tSTime, tETime time.Time) {
	var stInt, etInt int64
	var err error
	stString := c.Query("st")
	etString := c.Query("et")
	if stString == "" {
		tSTime = time.Now().AddDate(0, 0, -30)
	} else {
		stInt, err = strconv.ParseInt(stString, 10, 64)
		if err != nil {
			tSTime = time.Now().AddDate(0, 0, -30)
		} else {
			tSTime = time.Unix(stInt, 0)
		}
	}
	if etString == "" {
		tETime = time.Now()
	} else {
		etInt, err = strconv.ParseInt(etString, 10, 64)
		if err != nil {
			tETime = time.Now()
		} else {
			tETime = time.Unix(etInt, 0)
		}
	}
	return
}

func ReadZip(filebytes []byte) (str string, err error) {
	hashTmp := Md5String(time.Now().String())
	f, err := os.Create(hashTmp)
	defer f.Close()
	_, err = f.Write(filebytes)
	if err != nil {
		panic(err)
	}
	fi, err := f.Stat()
	if err != nil {
		panic(err)
		return
	}
	data, err := zip.NewReader(f, fi.Size())
	fs := data.File
	for i := range fs {
		fdata, err := fs[i].Open()
		if err != nil {
			break
		}
		b, _ := ioutil.ReadAll(fdata)
		str += string(b)
		fdata.Close()
	}
	os.Remove(hashTmp)
	return
}

func BytesToMd5String(data []byte) string {
	md5Ctx := md5.New()
	md5Ctx.Write(data)
	return hex.EncodeToString(md5Ctx.Sum(nil))
}

func Md5String(data string) string {
	return BytesToMd5String([]byte(data))
}

func CamelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return strings.Replace(string(data), "_", "", -1)
}


func SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

func GenerateMapStruct(d map[string]interface{}, name string) (structRaw string) {
	structRaw = fmt.Sprintf("type %v struct { \n", CamelString(name))
	var branch = ""
	var fieldType = make(map[string]string, len(d))
	for k := range d {
		switch v := d[k].(type) {
		case []interface{}:
			fieldType[k] = "interface{}"
		case map[string]interface{}:
			branch += GenerateMapStruct(v, CamelString(k))
			fieldType[k] = CamelString(k)
		case bson.ObjectId:
			fieldType[k] = "bson.ObjectId"
		case time.Time:
			fieldType[k] = "time.Time"
		case error:
			fieldType[k] = "error"
		default:
			if reflect.ValueOf(v).IsValid() {
				fieldType[k] = reflect.ValueOf(v).Kind().String()
			}
		}
		if fieldType[k] != "" {
			structRaw += fmt.Sprintf("\t%-20v %-20v `bson:\"%v\"`\n", CamelString(k), fieldType[k], k)
		}
	}
	structRaw = "\n" + structRaw + "}" + branch
	return structRaw
}
