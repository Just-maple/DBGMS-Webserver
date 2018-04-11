package session

import (
	"github.com/gin-contrib/sessions"
	"gopkg.in/mgo.v2/bson"
	"github.com/gin-gonic/gin"
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

func Default(c *gin.Context) UserSession {
	return UserSession{sessions.Default(c)}
}
