package controller

import (
	"github.com/Godzizizilla/douyin-simple/database"
	"github.com/Godzizizilla/douyin-simple/module"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func MessageAction(c *gin.Context) {
	currentUserID := c.MustGet("userID").(uint)
	toUserID, _ := strconv.Atoi(c.Query("to_user_id"))
	actionType := module.MessageAction(c.Query("action_type"))
	content := c.Query("content")

	now := time.Now()
	createTime := now.UnixNano() / int64(time.Millisecond)

	var message = module.Message{
		ToUserID:   uint(toUserID),
		FromUserID: currentUserID,
		Content:    content,
		CreateTime: createTime,
		CreatedAt:  now,
	}
	if actionType == module.Send {
		if err := database.MessageAction(&message); err == nil {
			c.JSON(http.StatusOK, module.Response{
				StatusCode: 0,
				StatusMsg:  "发送成功",
			})
			return
		}
	}
	c.JSON(http.StatusOK, module.Response{
		StatusCode: 1,
		StatusMsg:  "发送失败",
	})
}

func MessageChat(c *gin.Context) {
	currentUserID := c.MustGet("userID").(uint)
	toUserID, _ := strconv.Atoi(c.Query("to_user_id"))
	preMsgTime, _ := strconv.ParseInt(c.Query("pre_msg_time"), 10, 64)

	// preMsgTime = 0 -> 加载所有记录
	// preMsgTime = x -> 加载在这之后的数据
	messages := database.GetMessageListByUserID(uint(toUserID), currentUserID, preMsgTime)
	c.JSON(http.StatusOK, module.ChatResponse{
		Response:    module.Response{StatusCode: 0},
		MessageList: *messages,
	})
}
