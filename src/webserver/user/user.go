package user

import (
	"gopkg.in/mgo.v2/bson"
	"time"
	"webserver/dbx"
)

type DefaultUser struct {
	Id       bson.ObjectId `bson:"_id"`
	NickName string        `bson:"nickname"`
	Password string        `bson:"pwd"`
	TCreate  time.Time     `bson:"t_create"`
	TProcess time.Time     `bson:"t_process"`
	Level    UserLevel     `bson:"lvl"`
}

const (
	FieldNickName    = "nickname"
	FieldPassword    = "pwd"
	FieldTimeCreate  = "t_create"
	FieldTimeProcess = "t_process"
	FieldLevel       = "lvl"
	FieldId          = "_id"
	
	AllPermissionLevel = UserLevel(10)
)

type UserDBCollection struct {
	*dbx.Collection
}

type UserLevel int

func NewDefaultUser(nickname, password string, level UserLevel) (*DefaultUser) {
	return &DefaultUser{
		Id:       bson.NewObjectId(),
		NickName: nickname,
		Password: password,
		TCreate:  time.Now(),
		Level:    level,
	}
}

func (user *DefaultUser) HaveAllPermission() bool {
	return user.Level == AllPermissionLevel
}

func (c *UserDBCollection) UserLogin(nickname, password string) (user DefaultUser, err error) {
	err = c.Find(bson.M{FieldNickName: nickname, FieldPassword: password}).One(&user)
	return
}

func (c *UserDBCollection) GetUserData(Id bson.ObjectId) (user DefaultUser, err error) {
	err = c.FindId(Id).One(&user)
	return
}

func (c *UserDBCollection) NewUser(user *DefaultUser) (err error) {
	err = c.Insert(user)
	return
}

func (c *UserDBCollection) SetUserLevel(Id bson.ObjectId, level UserLevel) (err error) {
	return c.UpdateId(Id, bson.M{"$set": bson.M{FieldLevel: level}})
}
