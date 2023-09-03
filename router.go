package main

import (
	"github.com/Godzizizilla/douyin-simple/controller"
	"github.com/Godzizizilla/douyin-simple/middlewares"
	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	// basic apis
	apiRouter.GET("/feed", controller.Feed)
	apiRouter.GET("/user/", middlewares.JWT, controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	apiRouter.POST("/publish/action/", middlewares.JWT, controller.Publish)
	apiRouter.GET("/publish/list/", middlewares.JWT, controller.PublishList)

	// extra apis - I
	apiRouter.POST("/favorite/action/", middlewares.JWT, controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", middlewares.JWT, controller.FavoriteList)
	apiRouter.POST("/comment/action/", middlewares.JWT, controller.CommentAction)
	apiRouter.GET("/comment/list/", controller.CommentList)

	// extra apis - II
	apiRouter.POST("/relation/action/", middlewares.JWT, controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", middlewares.JWT, controller.FollowList)
	apiRouter.GET("/relation/follower/list/", middlewares.JWT, controller.FollowerList)
	apiRouter.GET("/relation/friend/list/", middlewares.JWT, controller.FriendList)
	apiRouter.GET("/message/chat/", middlewares.JWT, controller.MessageChat)
	apiRouter.POST("/message/action/", middlewares.JWT, controller.MessageAction)
}
