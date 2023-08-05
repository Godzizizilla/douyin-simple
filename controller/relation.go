package controller

import (
	"github.com/Godzizizilla/douyin-simple/database"
	"github.com/Godzizizilla/douyin-simple/module"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func RelationAction(c *gin.Context) {
	currentUserID := c.MustGet("userID").(uint)
	toUserID, _ := strconv.Atoi(c.Query("to_user_id"))
	actionType := module.RelationAction(c.Query("action_type"))

	if actionType != module.Follow && actionType != module.UnFollow {
		c.JSON(http.StatusOK, module.Response{
			StatusCode: 1,
			StatusMsg:  "关注/取消关注 失败",
		})
		return
	}

	if currentUserID == uint(toUserID) {
		c.JSON(http.StatusOK, module.Response{
			StatusCode: 1,
			StatusMsg:  "你不能关注自己哦",
		})
		return
	}

	if err := database.RelationAction(uint(toUserID), currentUserID, actionType); err != nil {
		c.JSON(http.StatusOK, module.Response{
			StatusCode: 1,
			StatusMsg:  "关注/取消关注 失败",
		})
		return
	}
	c.JSON(http.StatusOK, module.Response{
		StatusCode: 0,
		StatusMsg:  "关注/取消关注 成功",
	})
}

func FollowList(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))
	followings := database.GetFollowingListByUserID(uint(userID))
	c.JSON(http.StatusOK, module.UserListResponse{
		Response: module.Response{StatusCode: 0},
		UserList: *followings,
	})
}

func FollowerList(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))
	currentUserID := c.MustGet("userID").(uint)
	followers := database.GetFollowerListByUserID(uint(userID), currentUserID)
	c.JSON(http.StatusOK, module.UserListResponse{
		Response: module.Response{StatusCode: 0},
		UserList: *followers,
	})
}

func FriendList(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))
	currentUserID := c.MustGet("userID").(uint)
	friends := database.GetFriendListByUserID(uint(userID), currentUserID)
	c.JSON(http.StatusOK, module.UserListResponse{
		Response: module.Response{StatusCode: 0},
		UserList: *friends,
	})
}
