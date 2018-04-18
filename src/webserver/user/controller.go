package user

import (
	"webserver/dbx"
	"webserver/errorx"
	"time"
	"gopkg.in/mgo.v2/bson"
	"webserver/utilsx"
	"webserver/handler/controller"
)

type Controller struct {
	controller.DefaultController
	collection *dbx.Collection
}

func (c *Controller) userLogin(nickname, password string) (user DefaultUser, err error) {
	password = Md5EncodePassword(password)
	err = c.collection.Find(bson.M{FieldNickName: nickname, FieldPassword: password}).One(&user)
	if errorx.IsErrorNotFound(err) {
		err = errorx.ErrAuthFailed
	}
	return
}

func (c *Controller) changeUserPassword(Id bson.ObjectId, oldPWD, newPWD string) (err error) {
	oldPWD = Md5EncodePassword(oldPWD)
	newPWD = Md5EncodePassword(newPWD)
	err = c.collection.Update(bson.M{
		FieldId:       Id,
		FieldPassword: oldPWD,
	}, bson.M{
		dbx.BsonSelectorSet: bson.M{FieldPassword: newPWD,
		}})
	if errorx.IsErrorNotFound(err) {
		err = errorx.ErrAuthFailed
	}
	return
}

func (c *Controller) etUserLevelById(Id bson.ObjectId) (level Level) {
	user, err := c.getUserById(Id)
	if err != nil {
		return
	}
	return user.getUserLevel()
}

func (c *Controller) removeUserById(Id bson.ObjectId) (err error) {
	return c.collection.RemoveId(Id)
}

func (c *Controller) getUserById(Id bson.ObjectId) (user DefaultUser, err error) {
	err = c.collection.FindId(Id).One(&user)
	return
}

func (c *Controller) newUserFromNicknameAndPwd(nickname, password string, level Level, superiorUserId bson.ObjectId) (err error) {
	if c.checkUserNickNameValid(nickname) {
		password = Md5EncodePassword(password)
		var user = newUserFromNicknameAndPwd(nickname, password, level, superiorUserId)
		user.TCreate = time.Now()
		err = c.insertUser(user)
	} else {
		err = errorx.ErrAuthFailed
	}
	return
}

func (c *Controller) checkUserNickNameValid(nickname string) (bool) {
	err := c.collection.Find(bson.M{FieldNickName: nickname}).One(nil)
	return errorx.IsErrorNotFound(err)
}

func (c *Controller) insertUser(user *DefaultUser) (err error) {
	err = c.collection.Insert(user)
	return
}

func (c *Controller) setUserLevel(Id bson.ObjectId, level Level) (err error) {
	return c.collection.UpdateId(Id, bson.M{dbx.BsonSelectorSet: bson.M{FieldLevel: level}})
}

func (c *Controller) updateUserLogin(userId bson.ObjectId, ip string) (err error) {
	return c.collection.UpdateId(userId, bson.M{dbx.BsonSelectorSet: bson.M{FieldIP: ip, FieldTimeProcess: time.Now()}})
}

func (c *Controller) getAllUsers() (users []DefaultUser, err error) {
	err = c.collection.FindAll(nil, &users)
	return
}

func Md5EncodePassword(password string) string {
	return utilsx.Md5String(password + SecretSalt)
}
