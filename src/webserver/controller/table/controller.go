package table

import (
	"webserver/permission"
	"webserver/args"
	"webserver/controller"
	"webserver/dbx"
	"access"
)

type TableController struct {
	controller.DefaultController
	PermissionConfig *permission.Config
	Path             string
	Collection       *dbx.Collection
	GetAccessConfig  func(args *args.APIArgs) permission.AccessConfig
}

func InitAdminTableController() (controller *TableController) {
	return InitTableController(access.GetAdminPermissionConfig())
}
func InitTableController(cfg *permission.PermissionConfig) (controller *TableController) {
	controller = &TableController{
		PermissionConfig: &permission.Config{
			TableType: cfg.TableType,
			FieldType: cfg.FieldType,
		},
		Collection: cfg.Collection,
	}
	return
}

func (c *TableController) UseCollection(collection *dbx.Collection) {
	c.Collection = collection
}

func (c *TableController) Init() {
	var err error
	if c.Collection != nil {
		err = c.InitAllConfigTableFromDatabaseCollection()
	} else {
		err = c.InitAllConfigTableFromFiles()
	}
	if err != nil {
		log.Fatal(err)
	}
	c.RegisterAPI()
}

func (c *TableController) SetPath(path string) {
	c.Path = path
}

func (c *TableController) SetPermissionConfig(pc *permission.Config) {
	c.PermissionConfig = pc
}

func (c *TableController) SetAccessConfig(config func(args *args.APIArgs) permission.AccessConfig) {
	c.GetAccessConfig = config
}

func (c *TableController) GetPermissionConfig() *permission.Config {
	return c.PermissionConfig
}

