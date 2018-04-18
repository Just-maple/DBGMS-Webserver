package user

import "net/http"

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

const (
	FieldNickName    = "nickname"
	FieldPassword    = "pwd"
	FieldTimeCreate  = "t_create"
	FieldTimeProcess = "t_process"
	FieldLevel       = "lvl"
	FieldIP          = "ip"
	FieldId          = "_id"
	
	AllPermissionLevel = Level(10)
	SecretSalt         = "User-Secret-Salt"
)
