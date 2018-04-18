package user

import (
	"webserver/handler"
	"gopkg.in/mgo.v2/bson"
)

func (c *Controller) checkValidUser(args *handler.APIArgs) bool {
	valid, _ := args.UserId()
	return valid
}

func (c *Controller) registerUserPwdLoginApi() {
	c.RegisterApi(MethodPost, ApiAddrLogin, func(args *handler.APIArgs) (ret interface{}, err error) {
		nickname := args.JsonKey(JsonKeyNickname).MustString()
		password := args.JsonKey(JsonKeyPassword).MustString()
		user, err := c.userLogin(nickname, password)
		if err == nil {
			args.SetUserId(user.Id)
			ret = true
		}
		return
	})
}

func (c *Controller) registerNewUserApi() {
	c.RegisterApi(MethodPost, ApiAddrUser, func(args *handler.APIArgs) (ret interface{}, err error) {
		_, userId := args.UserId()
		nickname := args.JsonKey(JsonKeyNickname).MustString()
		password := args.JsonKey(JsonKeyPassword).MustString()
		level := args.JsonKey(JsonKeyLevel).MustInt()
		err = c.newUserFromNicknameAndPwd(nickname, password, Level(level), bson.ObjectIdHex(userId))
		return
	}, c.checkValidUser)
}

func (c *Controller) registerSetUserLevelApi() {
	c.RegisterApi(MethodPost, ApiAddrUserLevel, func(args *handler.APIArgs) (ret interface{}, err error) {
		userId := args.JsonKey(JsonKeyId).MustString()
		level := args.JsonKey(JsonKeyLevel).MustInt()
		err = c.setUserLevel(bson.ObjectIdHex(userId), Level(level))
		return
	}, c.checkValidUser)
}

func (c *Controller) registerChangeUserPasswordApi() {
	c.RegisterApi(MethodPost, ApiAddrPassword, func(args *handler.APIArgs) (ret interface{}, err error) {
		_, userId := args.UserId()
		oldPassword := args.JsonKey(JsonKeyOldPassword).MustString()
		newPassword := args.JsonKey(JsonKeyNewPassword).MustString()
		err = c.changeUserPassword(bson.ObjectIdHex(userId), oldPassword, newPassword)
		return
	}, c.checkValidUser)
}

func (c *Controller) registerGetAllUsersApi() {
	c.RegisterApi(MethodGet, ApiAddrAllUsers, func(args *handler.APIArgs) (ret interface{}, err error) {
		ret, err = c.getAllUsers()
		return
	}, c.checkValidUser)
}

func (c *Controller) registerUserSessionLoginApi() {
	c.RegisterApi(MethodGet, ApiAddrLogin, func(args *handler.APIArgs) (ret interface{}, err error) {
		_, userId := args.UserId()
		user, err := c.getUserById(bson.ObjectIdHex(userId))
		if err != nil {
			return
		}
		c.updateUserLogin(user.Id, args.IP())
		user.Password = ""
		ret = user
		return
	}, c.checkValidUser)
}
