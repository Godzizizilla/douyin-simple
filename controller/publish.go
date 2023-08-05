package controller

import (
	"fmt"
	"github.com/Godzizizilla/douyin-simple/database"
	"github.com/Godzizizilla/douyin-simple/module"
	"github.com/Godzizizilla/douyin-simple/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

var basicUrl = "http://192.168.124.2:8080"

func Publish(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	title := c.PostForm("title")
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, module.Response{
			StatusCode: 1,
			StatusMsg:  "上传失败",
		})
		return
	}

	// 保存视频
	timestamp := time.Now().Unix()
	videoFileName := fmt.Sprintf("%d_%d_%s", userID, timestamp, filepath.Base(data.Filename))
	videoSavePath := filepath.Join("./public/", videoFileName)
	videoSaveAbsolutePath, _ := filepath.Abs(videoSavePath)
	if err := c.SaveUploadedFile(data, videoSavePath); err != nil {
		c.JSON(http.StatusOK, module.Response{
			StatusCode: 1,
			StatusMsg:  "发布失败",
		})
		return
	}

	// 生成视频封面 ffmpeg -i input_video.mp4 -vf "select='eq(pict_type,PICT_TYPE_I)'" -vframes 1 output_frame.jpg
	imageFileName := utils.ChangeExtension(videoFileName, ".jpg")
	imageSavePath := filepath.Join("./public/", imageFileName)
	imageSaveAbsolutePath, _ := filepath.Abs(imageSavePath)

	cmd := exec.Command("ffmpeg", "-i", videoSaveAbsolutePath, "-vf", "select='eq(pict_type,PICT_TYPE_I)'", "-vframes", "1", imageSaveAbsolutePath)
	cmd.Run()

	// 存入数据库
	if err := database.AddVideo(&module.Video{
		UserID:   userID,
		Title:    title,
		PlayUrl:  basicUrl + "/static/" + videoFileName,
		CoverUrl: basicUrl + "/static/" + imageFileName,
	}); err != nil {
		c.JSON(http.StatusOK, module.Response{
			StatusCode: 1,
			StatusMsg:  "发布失败",
		})
		return
	}

	c.JSON(http.StatusOK, module.Response{
		StatusCode: 0,
		StatusMsg:  "发布成功",
	})
}

func PublishList(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))
	currentUserID := c.MustGet("userID").(uint)

	user, err := database.GetUserInfoByID(uint(userID), currentUserID)
	if err != nil {
		c.JSON(http.StatusOK, module.Response{
			StatusCode: 1,
			StatusMsg:  "用户ID错误",
		})
		return
	}

	videos, err := database.GetPublishVideosByUserID(uint(userID))
	if err != nil {
		c.JSON(http.StatusOK, module.Response{
			StatusCode: 0,
			StatusMsg:  "用户未发布视频",
		})
		return
	}

	for _, video := range *videos {
		video.Author = *user
	}

	c.JSON(http.StatusOK, module.VideoListResponse{
		Response:  module.Response{StatusCode: 0},
		VideoList: *videos,
	})

}
