package controller

import (
	"douyin-proj/src/global/ErrNo"
	"douyin-proj/src/global/util"
	"douyin-proj/src/service"
	"douyin-proj/src/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Feed(c *gin.Context) {
	// 请求参数获取
	var feedRequest = types.FeedRequest{}
	if err := c.ShouldBind(&feedRequest); err != nil {
		c.JSON(http.StatusOK, types.PublishResponse{StatusCode: ErrNo.ParamInvalid, StatusMsg: err.Error()})
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
	uid := uint(0)

	if token != "" {
		_uid, err := util.VerifyToken(token)
		if err != nil {
			c.JSON(http.StatusOK, types.FeedResponse{Response: ErrNo.AuthFailedResp})
			return
		}
		uid = _uid
		isAuth = true
	}
	feedVideos, nextTime, err := service.GetFeedVideos(lastTime, isAuth, uid)
	if err != nil {
		c.JSON(http.StatusOK, types.FeedResponse{Response: ErrNo.UnknownErrorResp})
		return
	}
	// 返回请求
	c.JSON(http.StatusOK, types.FeedResponse{
		Response:  ErrNo.SuccessResp,
		VideoList: feedVideos,
		NextTime:  nextTime,
	})
	return
}
