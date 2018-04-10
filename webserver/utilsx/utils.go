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
	f, err := os.Create("tmp")
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
	os.Remove("tmp")
	return
}

func BytesToMd5String(data []byte) string {
	md5Ctx := md5.New()
	md5Ctx.Write(data)
	return hex.EncodeToString(md5Ctx.Sum(nil))
}
func TransDate(Year, Month, Day int) (date time.Time) {
	year := strconv.Itoa(Year)
	month := strconv.Itoa(Month)
	day := strconv.Itoa(Day)
	if len(month) < 2 {
		month = "0" + month
	}
	if len(day) < 2 {
		day = "0" + day
	}
	date, _ = time.Parse("2006-01-02", year+"-"+month+"-"+day)
	return
}
