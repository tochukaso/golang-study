package controller

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/tochukaso/golang-study/env"
	csrf "github.com/utrack/gin-csrf"
)

func RenderHTML(c *gin.Context, status int, page string, params map[string]interface{}) {
	params["sUserID"], _ = c.Get("UserID")
	params["sUserName"], _ = c.Get("UserName")
	params["sUserCode"], _ = c.Get("UserCode")
	params["_csrf"] = csrf.GetToken(c)
	params["baseURL"] = env.GetEnv().BaseURL
	log.Println("params", params)
	c.HTML(status, page, params)
}

func ResponseJSON(c *gin.Context, status int, params map[string]interface{}) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
	c.JSON(status, params)
}
