package middlewares

import (
	"github.com/Godzizizilla/douyin-simple/module"
	"github.com/Godzizizilla/douyin-simple/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func JWT(c *gin.Context) {
	// get
	tokenString := c.Query("token")
	// post
	if tokenString == "" {
		tokenString = c.PostForm("token")
	}

	// 未提供Token
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, module.Response{
			StatusCode: 1,
			StatusMsg:  "未认证",
		})
		c.Abort()
		return
	}

	// 鉴权
	userID, err := utils.AuthenticateToken(tokenString)

	// 鉴权失败
	if err != nil {
		c.JSON(http.StatusUnauthorized, module.Response{
			StatusCode: 1,
			StatusMsg:  "token鉴权失败",
		})
		c.Abort()
		return
	}

	// 鉴权成功
	c.Set("userID", userID)
	c.Next()
}
