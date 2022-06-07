package server

import (
	"douyin-proj/src/config"
	"douyin-proj/src/controller"
	"github.com/gin-gonic/gin"
)

func Run(httpServer *gin.Engine) {
	// 设置日志格式,默认已经使用了
	//httpServer.Use(gin.Logger())
	httpServer.Use(gin.Recovery())

	// 注册路由
	initRouter(httpServer)

	serverError := httpServer.Run(":" + config.ServerConfig.HTTP_PORT)

	if serverError != nil {
		panic("server error !" + serverError.Error())
	}
}

func initRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/upload", "./upload")
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	apiRouter := r.Group("/douyin")

	// basic apis
	apiRouter.GET("/feed/", controller.Feed)
	apiRouter.GET("/user/", controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	apiRouter.POST("/publish/action/", controller.Publish)
	apiRouter.GET("/publish/list/", controller.PublishList)

	// extra apis - I
	apiRouter.POST("/favorite/action/", controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", controller.FavoriteList)
	apiRouter.POST("/comment/action/", controller.CommentAction)
	apiRouter.GET("/comment/list/", controller.CommentList)

	// extra apis - II
	apiRouter.POST("/relation/action/", controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", controller.FollowList)
	apiRouter.GET("/relation/follower/list/", controller.FollowerList)
}
