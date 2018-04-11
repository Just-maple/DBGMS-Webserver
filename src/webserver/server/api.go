package server

import (
	"github.com/bitly/go-simplejson"
	"github.com/gin-gonic/gin"
	"net/http"
	"webserver/jsonx"
	"webserver/session"
)

type JsonAPIFuncRoute map[string]JsonAPIFunc

func NewJsonAPIFuncRoute() JsonAPIFuncRoute {
	return make(JsonAPIFuncRoute)
}

func (h JsonAPIFuncRoute) RegisterAPI(name string, function JsonAPIFunc) {
	h[name] = function
}

func RenderJson(c *gin.Context, ok bool, data interface{}, err error) {
	c.Header("Access-Control-Allow-Origin", "*")
	if ok == true {
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

type JsonAPIFunc func(*gin.Context, *jsonx.Json, *session.UserSession) (interface{}, error)

func (svr *WebServer) JsonAPI(c *gin.Context) {
	var ok bool
	var ret interface{}
	var err error
	var function JsonAPIFunc
	var exists bool
	userSession := svr.ApiHandlers.GetSession(c)
	apiName := c.Param("api")
	function, exists = svr.ApiHandlers.GetApiFunc(c.Request.Method, apiName)
	if exists == true {
		var jsonData = new(jsonx.Json)
		if c.Request.Method == http.MethodPost {
			jsonData.Json, err = simplejson.NewFromReader(c.Request.Body)
		} else {
			jsonData.Json = simplejson.New()
		}
		if err == nil {
			ret, err = function(c, jsonData, userSession)
			if err == nil {
				ok = true
			}
		}

	}
	if svr.ApiHandlers.CheckDataBaseConnection(err); err == nil {
		ret = svr.ApiHandlers.RenderPermission(c, userSession, ret)
	} else {
		log.Errorf("JsonAPI(%s) err = %v", apiName, err)
	}
	RenderJson(c, ok, ret, err)
}
