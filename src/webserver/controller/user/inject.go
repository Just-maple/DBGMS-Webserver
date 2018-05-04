package user

import (
	"webserver/dbx"
	"webserver/controller"
)

func InitController(Collection *dbx.Collection) (*Controller) {
	return &Controller{DefaultController: controller.NewDefaultController(Collection)}
}

func (c *Controller) Init() {
	c.initApi()
	c.initUser()
}

func (c *Controller) initApi() {
	c.registerUserPwdLoginApi(ApiAddrLogin)
	c.registerNewUserApi(ApiAddrUser)
	c.registerSetUserLevelApi(ApiAddrUserLevel)
	c.registerUserSessionLoginAuthApi(ApiAuthLogin)
	c.registerChangeUserPasswordApi(ApiAddrPassword)
	c.registerLogOutApi(ApiAddrLogout)
	c.registerGetAllUsersApi(ApiAddrAllUsers)
}

func (c *Controller) initUser() {
	count, err := c.Collection.Count()
	if err != nil {
		log.Fatal("Get User Data Error", err)
	}
	if count == 0 {
		log.Debug("Not Found User,Init Default User admin")
		err = c.newUserFromNicknameAndPwd("admin", "admin", 0, "")
		if err != nil {
			log.Fatal(err)
		}
	}
}
