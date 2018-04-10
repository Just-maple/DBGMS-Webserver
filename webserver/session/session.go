package session

import (
	"github.com/gin-contrib/sessions"
	"gopkg.in/mgo.v2/bson"
)

const (
	SessionKeyUserId = "userid"
)

type UserSession struct {
	sessions.Session
}

func (us *UserSession) SetUserId(userId string) {
	us.Set(SessionKeyUserId, userId)
}

func (us *UserSession) AuthUserSession() (bool, string) {
	userId := us.Get(SessionKeyUserId)
	if userId == nil {
		return false, ""
	}
	return bson.IsObjectIdHex(userId.(string)), userId.(string)
}
