package controller

import (
	"douyin-proj/src/config"
	"douyin-proj/src/server/middleware"
	"douyin-proj/src/service"
	"douyin-proj/src/types"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

// Feed 获取视频流
func Feed(c *gin.Context) {
	// 请求参数获取
	var feedRequest = types.FeedRequest{}
	if err := c.ShouldBind(&feedRequest); err != nil {
		log.Printf("反序列化获取视频流请求失败。err:%s\n", err)
		c.JSON(http.StatusOK, types.PublishResponse{StatusCode: config.ParamInvalid, StatusMsg: err.Error()})
		return
	}
	// 截至时间
	lastTime := feedRequest.LastTime
	if lastTime == 0 {
		lastTime = time.Now().Unix()
	}
	// 鉴权
	token := feedRequest.Token
	isAuth := false
	userId := uint(0)
	if token != "" {
		_uid, err := middleware.VerifyToken(token)
		if err != nil {
			log.Printf("登录失败。err:%s\n", err)
			c.JSON(http.StatusOK, types.FeedResponse{Response: config.AuthFailedResp})
			return
		}
		userId = _uid
		isAuth = true
	}
	// 获得视频流
	feedVideos, nextTime, err := service.GetFeedVideosWithRecommendation(lastTime, isAuth, userId)
	if err != nil {
		c.JSON(http.StatusOK, types.FeedResponse{Response: config.UnknownErrorResp})
		return
	}
	// 返回请求
	c.JSON(http.StatusOK, types.FeedResponse{
		Response:  config.SuccessResp,
		VideoList: feedVideos,
		NextTime:  nextTime,
	})
	log.Printf("根据截至时间获取视频流成功。lastTime:%d\n", lastTime)
	return
}
