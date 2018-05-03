package user

import (
	"webserver/dbx"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func InitController(collection *dbx.Collection) (controller *Controller) {
	return &Controller{collection: collection}
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
	count, err := c.collection.Count()
	if err != nil {
		log.Fatal("Get User Data Error", err)
	}
	if count == 0 {
		log.Debug("Not Found User,Init Default User admin")
		uid := bson.NewObjectId()
		err = c.collection.Insert(DefaultUser{
			Id:           bson.NewObjectId(),
			NickName:     "admin",
			Password:     "admin",
			TCreate:      time.Now(),
			SuperiorUser: uid,
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}
