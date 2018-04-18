package user

import (
	"webserver/dbx"
	"webserver/handler"
)

func InitController(collection *dbx.Collection) (controller *Controller) {
	return &Controller{collection: collection}
}

func (c *Controller) InjectHandler(handler *handler.DefaultApiHandler) () {
	c.handler = handler
	c.registerUserPwdLoginApi()
	c.registerNewUserApi()
	c.registerSetUserLevelApi()
	c.registerUserSessionLoginApi()
	c.registerChangeUserPasswordApi()
	c.registerGetAllUsersApi()
}
