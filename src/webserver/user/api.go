package user

import (
	"webserver/handler"
	"net/http"
	"webserver/errorx"
	"gopkg.in/mgo.v2/bson"
)

const (
	MethodGet  = http.MethodGet
	MethodPost = http.MethodPost
	
	ApiAddrLogin     = "login"
	ApiAddrUserLevel = "userLevel"
	ApiAddrPassword  = "password"
	ApiAddrUser      = "user"
	ApiAddrAllUsers  = "users"
	
	JsonKeyNickname    = "nickname"
	JsonKeyPassword    = "password"
	JsonKeyLevel       = "level"
	JsonKeyId          = "Id"
	JsonKeyOldPassword = "oldPWD"
	JsonKeyNewPassword = "newPWD"
)

func (c *Controller) registerApi(method, api string, function handler.DefaultAPIFunc, pm ...handler.PermissionAuth) {
	h := c.handler.GetApiHandlersFromMethod(method)
	h.RegisterDefaultAPI(api, function, pm...)
}

func (c *Controller) RegisterUserPwdLoginApi() {
	c.registerApi(MethodPost, ApiAddrLogin, func(args *handler.APIArgs) (ret interface{}, err error) {
		nickname := args.JsonKey(JsonKeyNickname).MustString()
		password := args.JsonKey(JsonKeyPassword).MustString()
		user, err := c.UserLogin(nickname, password)
		if err == nil {
			args.SetUserId(user.Id)
			ret = true
		}
		return
	})
}

func (c *Controller) RegisterNewUserApi() {
	c.registerApi(MethodPost, ApiAddrUser, func(args *handler.APIArgs) (ret interface{}, err error) {
		nickname := args.JsonKey(JsonKeyNickname).MustString()
		password := args.JsonKey(JsonKeyPassword).MustString()
		level := args.JsonKey(JsonKeyLevel).MustInt()
		err = c.NewUserFromNicknameAndPwd(nickname, password, Level(level))
		return
	})
}

func (c *Controller) RegisterSetUserLevelApi() {
	c.registerApi(MethodPost, ApiAddrUserLevel, func(args *handler.APIArgs) (ret interface{}, err error) {
		userId := args.JsonKey(JsonKeyId).MustString()
		level := args.JsonKey(JsonKeyLevel).MustInt()
		err = c.setUserLevel(bson.ObjectIdHex(userId), Level(level))
		return
	})
}

func (c *Controller) RegisterChangeUserPasswordApi() {
	c.registerApi(MethodPost, ApiAddrPassword, func(args *handler.APIArgs) (ret interface{}, err error) {
		valid, userId := args.UserId()
		if !valid {
			err = errorx.ErrAuthFailed
			return
		}
		oldPassword := args.JsonKey(JsonKeyOldPassword).MustString()
		newPassword := args.JsonKey(JsonKeyNewPassword).MustString()
		err = c.ChangeUserPassword(bson.ObjectIdHex(userId), oldPassword, newPassword)
		return
	})
}

func (c *Controller) RegisterGetAllUsersApi() {
	c.registerApi(MethodGet, ApiAddrAllUsers, func(args *handler.APIArgs) (ret interface{}, err error) {
		ret, err = c.GetAllUsers()
		return
	})
}

func (c *Controller) RegisterUserSessionLoginApi() {
	c.registerApi(MethodGet, ApiAddrLogin, func(args *handler.APIArgs) (ret interface{}, err error) {
		valid, userId := args.UserId()
		if !valid {
			err = errorx.ErrAuthFailed
			return
		}
		user, err := c.GetUserById(bson.ObjectIdHex(userId))
		user.Password = ""
		ret = user
		return
	})
}
