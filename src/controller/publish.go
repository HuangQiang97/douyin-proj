package controller

import (
	"douyin-proj/src/config"
	"douyin-proj/src/server/middleware"
	"douyin-proj/src/service"
	"douyin-proj/src/types"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// Publish 上传视频
func Publish(c *gin.Context) {
	// 请求参数获取
	var publishRequest = types.PublishRequest{}
	if err := c.ShouldBind(&publishRequest); err != nil {
		log.Printf("反序列化上传视频请求失败。err:%s\n", err)
		c.JSON(http.StatusOK, types.PublishResponse{StatusCode: config.ParamInvalid, StatusMsg: err.Error()})
		return
	}
	// 校验jwt token
	userId, err := middleware.VerifyToken(publishRequest.Token)
	if err != nil {
		log.Printf("登录失败,err:%s\n", err)
		c.JSON(http.StatusOK, types.UserResponse{
			Response: config.AuthFailedResp,
		})
		return
	}
	// 保存文件
	err = service.SaveVideo(publishRequest.Data, userId, publishRequest.Title)
	if err != nil {
		// 创建文件失败，原因可能是1.路径不存在2.权限不足3.打开文件数量超过上限4.磁盘空间不足
		c.JSON(http.StatusOK, config.VideoUploadFailedResp)
		return
	}
	c.JSON(http.StatusOK, config.SuccessResp)
	log.Printf("保存文件成功。uid:%d\n", userId)
}

// PublishList 获取用户发表过的视频
func PublishList(c *gin.Context) {
	var videoListRequest = types.VideoListRequest{}
	if err := c.ShouldBind(&videoListRequest); err != nil {
		log.Printf("反序列化获取视频列表请求失败。err:%s\n", err)
		c.JSON(http.StatusOK, types.PublishResponse{StatusCode: config.ParamInvalid, StatusMsg: err.Error()})
		return
	}
	// 校验jwt token
	userId, err := middleware.VerifyToken(videoListRequest.Token)
	if err != nil {
		log.Printf("登录失败，err:%s\n", err)
		c.JSON(http.StatusOK, types.VideoListResponse{
			Response: config.AuthFailedResp,
		})
		return
	}
	// 获得视频列表
	videoList, err := service.GetVideoList(uint(videoListRequest.UserId), userId)
	if err != nil {
		c.JSON(http.StatusOK, types.VideoListResponse{
			Response: types.Response{
				StatusCode: config.UnknownError,
				StatusMsg:  "query user's videoList failed!"},
		})
		return
	}
	c.JSON(http.StatusOK, types.VideoListResponse{
		Response:  config.SuccessResp,
		VideoList: videoList,
	})
	log.Printf("获得用户发布过视频列表成功。uid:%d\n", videoListRequest.UserId)
}
