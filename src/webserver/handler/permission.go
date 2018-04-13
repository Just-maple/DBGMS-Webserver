package handler

import (
	"github.com/gin-gonic/gin"
	"reflect"
	. "webserver/session"
	"webserver/permission"
)

func (h *DefaultApiHandler) GetSession(c *gin.Context) (us *UserSession) {
	s := Default(c)
	us = &UserSession{Session: s}
	return
}


func (h *DefaultApiHandler) AuthUserSession(us *UserSession) (access permission.AccessConfig) {
	isValidId, userId := us.AuthUserSession()
	if isValidId {
		return h.db.GetAccessConfig(userId)
	} else {
		return nil
	}
}

func (h *DefaultApiHandler) RenderPermission(c *gin.Context, session *UserSession, in interface{}) (out interface{}) {
	if reflect.ValueOf(in).Kind() == reflect.Slice {
		config, has := h.PermissionConfig.GetConfigTableFromContext(c)
		if has {
			access := h.AuthUserSession(session)
			out = config.InitTablePermission(in, access)
		} else {
			out = in
		}
	} else {
		out = in
	}
	return
}
