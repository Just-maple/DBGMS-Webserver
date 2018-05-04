package table

import (
	"os"
	"io/ioutil"
	pm "webserver/permission"
)

type TableData []struct {
	Name string `bson:"_id"`
	Data string `bson:"data"`
}

func (c *Controller) InitAllConfigTableFromDatabaseCollection() (err error) {
	var tableData TableData
	err = c.Collection.Find(nil).All(&tableData)
	if err != nil {
		return
	}
	c.PermissionConfig.TableMap = make(map[string]*pm.Table, len(tableData))
	for i := range tableData {
		err = c.PermissionConfig.InitTableConfig([]byte(tableData[i].Data), tableData[i].Name)
		if err != nil {
			log.Fatal(err)
		}
	}
	return
}

func (c *Controller) InitAllConfigTableFromFiles() (err error) {
	if !IsExist(c.Path) {
		err = os.Mkdir(c.Path, 0700)
		if err != nil {
			return
		}
	}
	tableFiles, err := ioutil.ReadDir(c.Path)
	if err != nil {
		return
	}
	c.PermissionConfig.TableMap = make(map[string]*pm.Table, len(tableFiles))
	for i := range tableFiles {
		err = c.initTableConfigFromFileInfo(&(tableFiles[i]))
		if err != nil {
			log.Fatal(err)
		}
	}
	return
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
