package controller

import (
	"github.com/Godzizizilla/douyin-simple/database"
	"github.com/Godzizizilla/douyin-simple/module"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func Feed(c *gin.Context) {
	var timestamp time.Time
	unixTimeString := c.Query("latest_time")
	if unixTimeString == "" {
		timestamp = time.Now()
	}
	unixTime, _ := strconv.Atoi(unixTimeString)
	timestamp = time.Unix(int64(unixTime), 0)

	videos, err := database.GetVideosBeforeTimestamp(timestamp)
	if err != nil {
		c.JSON(http.StatusOK, module.Response{
			StatusCode: 1,
			StatusMsg:  "没有更早的视频了",
		})
	}

	c.JSON(http.StatusOK, module.FeedResponse{
		Response:  module.Response{StatusCode: 0},
		VideoList: *videos,
		NextTime:  time.Now().Unix(),
	})
}
