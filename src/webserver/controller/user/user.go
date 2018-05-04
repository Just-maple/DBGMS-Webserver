package user

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type DefaultUser struct {
	Id           bson.ObjectId `bson:"_id"`
	NickName     string        `bson:"nickname"`
	Password     string        `bson:"pwd"`
	TCreate      time.Time     `bson:"t_create"`
	TProcess     time.Time     `bson:"t_process"`
	Level        Level         `bson:"lvl"`
	SuperiorUser bson.ObjectId `bson:"sp_user"`
	IP           string        `bson:"ip"`
}

type Level = int

func newUserFromNicknameAndPwd(nickname, hashPassword string, level Level, SuperiorUserId bson.ObjectId) *DefaultUser {
	return &DefaultUser{
		Id:           bson.NewObjectId(),
		NickName:     nickname,
		Password:     hashPassword,
		TCreate:      time.Now(),
		Level:        level,
		SuperiorUser: SuperiorUserId,
	}
}

func (user *DefaultUser) getUserLevel() Level {
	return user.Level
}

func (user *DefaultUser) haveAllPermission() bool {
	return user.getUserLevel() >= AllPermissionLevel
}
