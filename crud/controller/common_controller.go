package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func RenderHTML(c *gin.Context, status int, page string, params map[string]interface{}) {
	params["sUserID"], _ = c.Get("UserID")
	params["sUserName"], _ = c.Get("UserName")
	params["sUserCode"], _ = c.Get("UserCode")
	params["_csrf"] = csrf.GetToken(c)
	fmt.Println("params", params)
	c.HTML(status, page, params)
}
