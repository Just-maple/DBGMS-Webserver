package user

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type DefaultUser struct {
	Id       bson.ObjectId `bson:"_id"`
	NickName string        `bson:"nickname"`
	Password string        `bson:"pwd"`
	TCreate  time.Time     `bson:"t_create"`
	TProcess time.Time     `bson:"t_process"`
	Level    Level         `bson:"lvl"`
}

const (
	FieldNickName    = "nickname"
	FieldPassword    = "pwd"
	FieldTimeCreate  = "t_create"
	FieldTimeProcess = "t_process"
	FieldLevel       = "lvl"
	FieldId          = "_id"
	
	AllPermissionLevel = Level(10)
)

type Level int

const SecretSalt = "User-Secret-Salt"

func NewUserFromNicknameAndPwd(nickname, hashPassword string, level Level) (*DefaultUser) {
	return &DefaultUser{
		Id:       bson.NewObjectId(),
		NickName: nickname,
		Password: hashPassword,
		TCreate:  time.Now(),
		Level:    level,
	}
}

func (user *DefaultUser) GetUserLevel() Level {
	return user.Level
}

func (user *DefaultUser) HaveAllPermission() bool {
	return user.GetUserLevel() >= AllPermissionLevel
}
