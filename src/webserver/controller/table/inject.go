package table

import (
	"webserver/permission"
	"webserver/args"
	"webserver/controller"
	"webserver/dbx"
	"access"
)

type Controller struct {
	*controller.DefaultController
	PermissionConfig  *permission.Config
	Path              string
	OauthAccessConfig func(args *args.APIArgs) permission.AccessConfig
}

func InitAdminTableController(tablePath string, collection *dbx.Collection) (controller *Controller) {
	return InitTableController(access.GetAdminPermissionConfig(), tablePath, collection)
}

func InitAdminTableControllerByCollection(collection *dbx.Collection) (controller *Controller) {
	return InitTableController(access.GetAdminPermissionConfig(), "", collection)
}

func InitAdminTableControllerByPath(tablePath string) (controller *Controller) {
	return InitTableController(access.GetAdminPermissionConfig(), tablePath, nil)
}

func InitTableController(cfg *permission.PermissionConfig, tablePath string, collection *dbx.Collection) (tableController *Controller) {
	tableController = &Controller{
		PermissionConfig: &permission.Config{
			TableType: cfg.TableType,
			FieldType: cfg.FieldType,
		},
		Path:              tablePath,
		DefaultController: controller.NewDefaultController(collection),
	}
	return
}

func (c *Controller) Init() {
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

func (c *Controller) SetAccessConfig(config func(args *args.APIArgs) permission.AccessConfig) {
	c.OauthAccessConfig = config
}

func (c *Controller) GetPermissionConfig() *permission.Config {
	return c.PermissionConfig
}
