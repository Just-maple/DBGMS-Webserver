package handler

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"reflect"
	. "webserver/session"
)

func (h *DefaultApiHandler) GetSession(c *gin.Context) (us *UserSession) {
	s := sessions.Default(c)
	us = &UserSession{Session: s}
	return
}

func (h *DefaultApiHandler) AuthAdminUserSession(us *UserSession) bool {
	isValidId, userId := us.AuthUserSession()
	return isValidId && h.db.AuthAdminUser(userId)
}

func (h *DefaultApiHandler) AuthSuperAdminUserSession(us *UserSession) (isAdmin, isSuper bool) {
	isValidId, userId := us.AuthUserSession()
	if isValidId {
		return h.db.AuthSuperAdminUser(userId)
	} else {
		return
	}
}

func (h *DefaultApiHandler) RenderPermission(c *gin.Context, session *UserSession, in interface{}) (out interface{}) {
	if reflect.ValueOf(in).Kind() == reflect.Slice {
		config, has := h.TableConfig.GetConfigTableFromContext(c)
		if has {
			isAdmin, isSuper := h.AuthSuperAdminUserSession(session)
			out = config.InitTablePermission(in, isAdmin, isSuper)
		} else {
			out = in
		}
	} else {
		out = in
	}
	return
}
