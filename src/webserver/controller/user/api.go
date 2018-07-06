package user

import (
	"gopkg.in/mgo.v2/bson"
	"logger"
	"webserver/args"
	"webserver/errorx"
)

var log = logger.Log

func (c *Controller) CompareUserLevel(args *args.APIArgs) bool {
	valid, userId := args.UserId()
	operatedId := args.JsonKeyId()
	return !valid || bson.IsObjectIdHex(operatedId) || c.getUserLevelById(bson.ObjectIdHex(userId)) > c.getUserLevelById(bson.ObjectIdHex(operatedId))
}

func (c *Controller) checkValidUser(args *args.APIArgs) bool {
	valid, _ := args.UserId()
	return valid
}

func (c *Controller) registerUserPwdLoginApi(ApiAddrLogin string) {
	c.RegisterPostApi(ApiAddrLogin, func(args *args.APIArgs) (ret interface{}, err error) {
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

func (c *Controller) registerLogOutApi(api string) {
	c.RegisterGetApi(api, func(args *args.APIArgs) (ret interface{}, err error) {
		args.ClearSession()
		return
	}, c.checkValidUser)
}

func (c *Controller) registerNewUserApi(ApiAddrUser string) {
	c.RegisterPostApi(ApiAddrUser, func(args *args.APIArgs) (ret interface{}, err error) {
		_, userId := args.UserId()
		nickname := args.JsonKey(JsonKeyNickname).MustString()
		password := args.JsonKey(JsonKeyPassword).MustString()
		level := args.JsonKey(JsonKeyLevel).MustInt()
		err = c.newUserFromNicknameAndPwd(nickname, password, Level(level), userId)
		return
	}, c.checkValidUser)
}

func (c *Controller) registerSetUserLevelApi(ApiAddrUserLevel string) {
	c.RegisterPostApi(ApiAddrUserLevel, func(args *args.APIArgs) (ret interface{}, err error) {
		userId := args.JsonKeyId()
		level := args.JsonKey(JsonKeyLevel).MustInt()
		err = c.setUserLevel(bson.ObjectIdHex(userId), Level(level))
		return
	}, c.CompareUserLevel)
}

func (c *Controller) registerChangeUserPasswordApi(ApiAddrPassword string) {
	c.RegisterPostApi(ApiAddrPassword, func(args *args.APIArgs) (ret interface{}, err error) {
		_, userId := args.UserId()
		oldPassword := args.JsonKey(JsonKeyOldPassword).MustString()
		newPassword := args.JsonKey(JsonKeyNewPassword).MustString()
		err = c.changeUserPassword(bson.ObjectIdHex(userId), oldPassword, newPassword)
		return
	}, c.checkValidUser)
}

func (c *Controller) registerGetAllUsersApi(ApiAddrAllUsers string) {
	c.RegisterGetApi(ApiAddrAllUsers, func(args *args.APIArgs) (ret interface{}, err error) {
		ret, err = c.getAllUsers()
		return
	}, c.checkValidUser)
}

func (c *Controller) registerUserSessionLoginAuthApi(ApiAuthLogin string) {
	c.RegisterGetApi(ApiAuthLogin, func(args *args.APIArgs) (ret interface{}, err error) {
		valid, userId := args.UserId()
		if !valid {
			err = errorx.ErrAuthFailed
			return
		}
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
