package user

import (
	"webserver/dbx"
	"webserver/errorx"
	"time"
	"gopkg.in/mgo.v2/bson"
	"webserver/utilsx"
	"webserver/handler"
)

type Controller struct {
	collection *dbx.Collection
	handler    *handler.DefaultApiHandler
}

func (c *Controller) UserLogin(nickname, password string) (user DefaultUser, err error) {
	password = Md5EncodePassword(password)
	err = c.collection.Find(bson.M{FieldNickName: nickname, FieldPassword: password}).One(&user)
	if errorx.IsErrorNotFound(err) {
		err = errorx.ErrAuthFailed
	}
	return
}

func (c *Controller) ChangeUserPassword(Id bson.ObjectId, oldPWD, newPWD string) (err error) {
	oldPWD = Md5EncodePassword(oldPWD)
	newPWD = Md5EncodePassword(newPWD)
	err = c.collection.Update(bson.M{FieldId: Id, FieldPassword: oldPWD}, bson.M{"$set": bson.M{FieldPassword: newPWD}})
	if errorx.IsErrorNotFound(err) {
		err = errorx.ErrAuthFailed
	}
	return
}

func (c *Controller) GetUserLevelById(Id bson.ObjectId) (level Level) {
	user, err := c.GetUserById(Id)
	if err != nil {
		return
	}
	return user.GetUserLevel()
}

func (c *Controller) RemoveUserById(Id bson.ObjectId) (err error) {
	return c.collection.RemoveId(Id)
}

func (c *Controller) GetUserById(Id bson.ObjectId) (user DefaultUser, err error) {
	err = c.collection.FindId(Id).One(&user)
	return
}

func (c *Controller) NewUserFromNicknameAndPwd(nickname, password string, level Level) (err error) {
	if c.checkUserNickNameValid(nickname) {
		password = Md5EncodePassword(password)
		var user = NewUserFromNicknameAndPwd(nickname, password, level)
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
	return c.collection.UpdateId(Id, bson.M{"$set": bson.M{FieldLevel: level}})
}

func (c *Controller) GetAllUsers() (users []DefaultUser, err error) {
	err = c.collection.FindAll(nil, &users)
	return
}

func Md5EncodePassword(password string) string {
	return utilsx.Md5String(password + SecretSalt)
}
