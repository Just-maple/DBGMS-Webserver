package user

import (
	"webserver/dbx"
)

func InitController(collection *dbx.Collection) (controller *Controller) {
	return &Controller{collection: collection}
}

func (c *Controller) Init() {
	c.registerUserPwdLoginApi(ApiAddrLogin)
	c.registerNewUserApi(ApiAddrUser)
	c.registerSetUserLevelApi(ApiAddrUserLevel)
	c.registerUserSessionLoginAuthApi(ApiAuthLogin)
	c.registerChangeUserPasswordApi(ApiAddrPassword)
	c.registerLogOutApi(ApiAddrLogout)
	c.registerGetAllUsersApi(ApiAddrAllUsers)
}
