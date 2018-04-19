package user

import (
	"webserver/dbx"
)

func InitController(collection *dbx.Collection) (controller *Controller) {
	return &Controller{collection: collection}
}

func (c *Controller) Init() {
	c.registerUserPwdLoginApi()
	c.registerNewUserApi()
	c.registerSetUserLevelApi()
	c.registerUserSessionLoginApi()
	c.registerChangeUserPasswordApi()
	c.registerGetAllUsersApi()
}
