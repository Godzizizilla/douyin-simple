package controller

import (
	"github.com/Godzizizilla/douyin-simple/database"
	"github.com/Godzizizilla/douyin-simple/module"
	"github.com/Godzizizilla/douyin-simple/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	_, _, err := database.GetUserByUsername(username)
	if err == nil {
		c.JSON(http.StatusOK, module.Response{
			StatusCode: 1,
			StatusMsg:  "用户已存在或用户名重复",
		})
		return
	}

	userID, err := database.AddUser(username, password)
	if err != nil {
		c.JSON(http.StatusOK, module.Response{
			StatusCode: 1,
			StatusMsg:  "注册失败",
		})
		return
	}

	token, _ := utils.GenerateToken(userID)
	c.JSON(http.StatusOK, module.UserResponse{
		Response: module.Response{StatusCode: 0},
		UserId:   userID,
		Token:    token,
	})
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	// 用户是否存在, 存在返回PasswordHash
	userID, encryptedPassword, err := database.GetUserByUsername(username)
	if err != nil {
		c.JSON(http.StatusOK, module.Response{
			StatusCode: 1,
			StatusMsg:  "用户不存在",
		})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(password)); err != nil {
		c.JSON(http.StatusOK, module.Response{
			StatusCode: 1,
			StatusMsg:  "密码错误",
		})
		return
	}

	token, _ := utils.GenerateToken(userID)
	c.JSON(http.StatusOK, module.UserResponse{
		Response: module.Response{StatusCode: 0},
		UserId:   userID,
		Token:    token,
	})
}

func UserInfo(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))

	// 获取用户信息
	userInfo, err := database.GetUserInfoByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusOK, module.Response{
			StatusCode: 1,
			StatusMsg:  "用户不存在",
		})
		return
	}
	c.JSON(http.StatusOK, module.UserInfoResponse{
		Response: module.Response{StatusCode: 0},
		User:     *userInfo,
	})
}
