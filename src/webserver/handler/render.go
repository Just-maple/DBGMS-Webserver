package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *DefaultApiHandler) JsonAPI(c *gin.Context) {
	var ok bool
	var ret interface{}
	var err error
	var userSession = h.GetSession(c)
	var apiName = c.Param("api")
	function, exists := h.GetApiFunc(c.Request.Method, apiName)
	if exists {
		ret, err = function.Run(c, userSession)
		if h.CheckDataBaseConnection(err); err == nil {
			ret = h.RenderPermission(c, userSession, ret)
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

