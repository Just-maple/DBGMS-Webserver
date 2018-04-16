package session

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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
	userId, valid := us.Get(SessionKeyUserId).(string)
	if !valid {
		return false, ""
	}
	return bson.IsObjectIdHex(userId), userId
}

func Default(c *gin.Context) UserSession {
	return UserSession{sessions.Default(c)}
}
