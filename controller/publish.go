package controller

import (
	"fmt"
	"github.com/Godzizizilla/douyin-simple/database"
	"github.com/Godzizizilla/douyin-simple/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

var basicUrl = "http://172.22.22.94:8080/"

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	// TODO JWT认证
	token := c.PostForm("token")
	claims, err := utils.AuthenticateToken(token)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "用户不存在",
		})
		return
	}

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	title := c.PostForm("title")

	timestamp := time.Now().Unix()

	filename := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d_%d_%s", claims.ID, timestamp, filename)
	saveFile := filepath.Join("./public/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	// 生成视频封面
	// ffmpeg -i input_video.mp4 -vf "select='eq(pict_type,PICT_TYPE_I)'" -vframes 1 output_frame.jpg
	outputImagePath := utils.ChangeExtension(saveFile, ".jpg")
	saveFileAbsolutePath, _ := filepath.Abs(saveFile)
	outputImagePathAbsolutePath, _ := filepath.Abs(outputImagePath)

	cmd := exec.Command("ffmpeg", "-i", saveFileAbsolutePath, "-vf", "select='eq(pict_type,PICT_TYPE_I)'", "-vframes", "1", outputImagePathAbsolutePath)
	output, err2 := cmd.CombinedOutput()
	if err2 != nil {
		fmt.Println("Error:", err2)
		fmt.Println("Output:", string(output))
		return
	}

	// TODO 将视频url存入数据库
	if err := database.AddVideo(&database.Video{
		UserID:        claims.ID,
		Title:         title,
		PlayUrl:       saveFile,
		CoverUrl:      outputImagePath,
		FavoriteCount: 0,
		CommentCount:  0,
	}); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	token := c.Query("token")
	_, err := utils.AuthenticateToken(token)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "用户不存在",
		})
		return
	}

	userID, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "用户ID错误",
		})
		return
	}

	var videos []database.Video
	database.GetVideosByUserID(userID, &videos)

	userInfo, err := database.GetUserInfoByID(userID)
	user := User{
		Id:            userInfo.ID,
		Name:          userInfo.Name,
		FollowCount:   int64(userInfo.FollowCount),
		FollowerCount: int64(userInfo.FollowerCount),
		IsFollow:      false,
	}
	var videoList []Video
	for _, video := range videos {
		videoList = append(videoList, Video{
			Id:            video.ID,
			Author:        user,
			PlayUrl:       basicUrl + video.PlayUrl,
			CoverUrl:      basicUrl + video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			// TODO
			IsFavorite: false,
		})
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videoList,
	})
}
