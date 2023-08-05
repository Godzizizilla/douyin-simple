package controller

import (
	"github.com/Godzizizilla/douyin-simple/database"
	"github.com/Godzizizilla/douyin-simple/module"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	videoID, _ := strconv.Atoi(c.Query("video_id"))
	actionType := c.Query("action_type")

	if actionType != "1" && actionType != "2" {
		c.JSON(http.StatusOK, module.Response{
			StatusCode: 1,
			StatusMsg:  "点赞/取消点赞 失败",
		})
	}

	if err := database.FavoriteAction(userID, uint(videoID), module.FavoriteAction(actionType)); err != nil {
		c.JSON(http.StatusOK, module.Response{
			StatusCode: 1,
			StatusMsg:  "点赞/取消点赞 失败",
		})
	}
	c.JSON(http.StatusOK, module.Response{
		StatusCode: 0,
		StatusMsg:  "成功",
	})
}

func FavoriteList(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))
	currentUserID := c.MustGet("userID").(uint)

	videos := database.GetFavoriteVideosByUserID(uint(userID), currentUserID)
	c.JSON(http.StatusOK, module.VideoListResponse{
		Response:  module.Response{StatusCode: 0},
		VideoList: *videos,
	})

}
