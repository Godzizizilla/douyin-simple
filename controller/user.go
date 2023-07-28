package controller

import (
	"github.com/Godzizizilla/douyin-simple/database"
	"github.com/Godzizizilla/douyin-simple/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

var userIdSequence = int64(1)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	if exist, _ := database.IsUserExists(username); exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		// 对密码进行加密
		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			panic("Error hashing password:")
		}
		// 创建用户数据
		user := database.User{Name: username, PasswordHash: hashedPassword}
		// 添加用户
		userID, err := database.AddUser(&user)
		// 添加失败
		if err != nil {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: "register failed"},
			})
		}
		// 添加成功, 生成token
		token, err := utils.GenerateToken(userID)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   userID,
			Token:    token,
		})
	}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	// 用户是否存在, 存在返回PasswordHash
	exist, hashedPassword := database.IsUserExists(username)
	if !exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
	// 用户密码是否正确
	if err := utils.VerifyPassword(password, hashedPassword); err == nil {
		userID, err := database.GetUserIdByName(username)
		if err != nil {
			println(err)
			return
		}
		token, err := utils.GenerateToken(userID)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   userID,
			Token:    token,
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Password error"},
		})
	}
}

func UserInfo(c *gin.Context) {
	token := c.Query("token")
	claims, err := utils.AuthenticateToken(token)
	if err == nil {
		// 获取用户信息
		userInfo, err := database.GetUserInfoByID(claims.ID)
		if err == nil {
			user := User{
				Id:            userInfo.ID,
				Name:          userInfo.Name,
				FollowCount:   int64(userInfo.FollowCount),
				FollowerCount: int64(userInfo.FollowerCount),
				IsFollow:      false,
			}
			c.JSON(http.StatusOK, UserResponse{
				Response: Response{StatusCode: 0},
				User:     user,
			})
		} else {
			c.JSON(http.StatusOK, UserResponse{
				Response: Response{StatusCode: 1, StatusMsg: "get User Info failed"},
			})
		}
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}
