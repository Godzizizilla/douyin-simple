package controller

import (
	"github.com/Godzizizilla/douyin-simple/database"
	"github.com/Godzizizilla/douyin-simple/module"
	"github.com/Godzizizilla/douyin-simple/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func Feed(c *gin.Context) {
	// TODO 返回最早视频的时间戳

	tokenString := c.Query("token")
	// 不需要处理error, 因为有没有token都可
	userID, _ := utils.AuthenticateToken(tokenString)

	var timestamp time.Time
	unixTimeString := c.Query("latest_time")
	unixTime, _ := strconv.Atoi(unixTimeString)
	timestamp = time.Unix(int64(unixTime), 0)

	if unixTimeString == "" || int64(unixTime) > time.Now().Unix() {
		timestamp = time.Now()
	}

	videos, err := database.GetVideosBeforeTimestamp(timestamp, userID)
	if err != nil || len(*videos) == 0 {
		c.JSON(http.StatusOK, module.Response{
			StatusCode: 1,
			StatusMsg:  "没有更早的视频了",
		})
		return
	}

	lastVideo := &(*videos)[len(*videos)-1]

	c.JSON(http.StatusOK, module.FeedResponse{
		Response:  module.Response{StatusCode: 0},
		VideoList: *videos,
		NextTime:  lastVideo.CreatedAt.Unix(),
	})
}
