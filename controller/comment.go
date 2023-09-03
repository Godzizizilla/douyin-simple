package controller

import (
	"fmt"
	"github.com/Godzizizilla/douyin-simple/database"
	"github.com/Godzizizilla/douyin-simple/module"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func CommentAction(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	videoID, _ := strconv.Atoi(c.Query("video_id"))
	actionType := module.CommentAction(c.Query("action_type"))

	if actionType != module.PublishComment && actionType != module.DeleteComment {
		c.JSON(http.StatusOK, module.Response{
			StatusCode: 1,
			StatusMsg:  "发布/删除评论 失败",
		})
	}

	var comment module.Comment

	if actionType == module.PublishComment {
		now := time.Now()
		month := now.Month()
		day := now.Day()
		dateString := fmt.Sprintf("%02d-%02d", month, day)

		comment.Content = c.Query("comment_text")
		comment.UserID = userID
		comment.VideoID = uint(videoID)
		comment.CreateDate = dateString
		if err := database.CommentAction(&comment, actionType); err == nil {
			c.JSON(http.StatusOK, module.CommentActionResponse{
				Response: module.Response{StatusCode: 0},
				Comment:  comment,
			})
			return
		}

	} else if actionType == module.DeleteComment {
		commentID, _ := strconv.Atoi(c.Query("comment_id"))
		comment.ID = uint(commentID)
		comment.VideoID = uint(videoID)
		fmt.Println(commentID, videoID)
		if err := database.CommentAction(&comment, actionType); err == nil {
			c.JSON(http.StatusOK, module.Response{
				StatusCode: 0,
				StatusMsg:  "删除评论成功",
			})
			return
		}
	}
	c.JSON(http.StatusOK, module.Response{
		StatusCode: 1,
		StatusMsg:  "发布/删除评论 失败",
	})

}

func CommentList(c *gin.Context) {
	// userID := c.MustGet("userID").(uint)
	videoID, _ := strconv.Atoi(c.Query("video_id"))

	comments := database.GetCommentsByVideoID(uint(videoID))
	c.JSON(http.StatusOK, module.CommentListResponse{
		Response:    module.Response{StatusCode: 0},
		CommentList: *comments,
	})
}
