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
	c.RegisterUserPwdLoginApi()
	c.RegisterNewUserApi()
	c.RegisterSetUserLevelApi()
	c.RegisterUserSessionLoginApi()
	c.RegisterChangeUserPasswordApi()
	c.RegisterGetAllUsersApi()
}
