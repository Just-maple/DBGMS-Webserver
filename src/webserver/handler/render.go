package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"webserver/args/jsonx"
	"webserver/args/session"
	"webserver/args"
)

func (h *DefaultApiHandler) JsonAPI(c *gin.Context) {
	var ok bool
	var ret interface{}
	var err error
	var userSession = &session.UserSession{Session: session.Default(c)}
	var apiName = c.Param("api")
	function, exists := h.getApiFunc(c.Request.Method, apiName)
	if exists {
		var jsonData *jsonx.Json
		var arg *args.APIArgs
		switch c.Request.Method {
		case http.MethodGet:
			jsonData = jsonx.New()
		default:
			jsonData, err = jsonx.NewFromReader(c.Request.Body)
		}
		if err == nil {
			arg = args.New(c, jsonData, userSession)
			ret, err = function.Run(arg)
		}
		if h.CheckDataBaseConnection(err); err == nil {
			ret = h.RenderPermission(arg, ret)
			ok = true
		} else {
			log.Errorf("JsonAPI(%s) err = %v", apiName, err)
		}
	}
	RenderJson(c, ok, ret, err)
}

func RenderJson(c *gin.Context, ok bool, data interface{}, err error) {
	c.Header("Access-Control-Allow-Origin", "*")
	if ok {
		c.JSON(http.StatusOK, map[string]interface{}{
			"ok":   ok,
			"data": data,
		})
	} else {
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, map[string]interface{}{
				"ok":   ok,
				"data": data,
				"err":  err.Error(),
			})
			
		} else {
			c.IndentedJSON(http.StatusInternalServerError, map[string]interface{}{
				"ok":   ok,
				"data": data,
			})
		}
	}
}
