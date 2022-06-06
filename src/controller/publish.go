package controller

import (
	"douyin-proj/src/global/ErrNo"
	"douyin-proj/src/global/util"
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
		c.JSON(http.StatusOK, types.PublishResponse{StatusCode: ErrNo.ParamInvalid, StatusMsg: err.Error()})
		return
	}

	// 校验jwt token
	uId, err := util.VerifyToken(publishRequest.Token)
	if err != nil {
		log.Printf("登录失败，err:%s\n", err)
		c.JSON(http.StatusOK, types.UserResponse{
			Response: ErrNo.AuthFailedResp,
		})
		return
	}

	// 保存文件
	err = service.SaveVideo(publishRequest.Data, uId, publishRequest.Title)
	if err != nil {
		// 创建文件失败，原因可能是1.路径不存在2.权限不足3.打开文件数量超过上限4.磁盘空间不足
		c.JSON(http.StatusOK, ErrNo.VideoUploadFailedResp)
		return
	}
	c.JSON(http.StatusOK, ErrNo.SuccessResp)
	log.Printf("保存文件成功。uid:%d\n", uId)

}

// PublishList 获取用户发表过的视频
func PublishList(c *gin.Context) {
	var videoListRequest = types.VideoListRequest{}
	if err := c.ShouldBind(&videoListRequest); err != nil {
		log.Printf("反序列化获取视频列表请求失败。err:%s\n", err)
		c.JSON(http.StatusOK, types.PublishResponse{StatusCode: ErrNo.ParamInvalid, StatusMsg: err.Error()})
		return
	}
	// 校验jwt token
	uId, err := util.VerifyToken(videoListRequest.Token)
	if err != nil {
		log.Printf("登录失败，err:%s\n", err)
		c.JSON(http.StatusOK, types.VideoListResponse{
			Response: ErrNo.AuthFailedResp,
		})
		return
	}

	// 获得视频列表
	videoList, err := service.GetVideoList(uint(videoListRequest.UserId), uId)
	if err != nil {
		c.JSON(http.StatusOK, types.VideoListResponse{
			Response: types.Response{
				StatusCode: ErrNo.UnknownError,
				StatusMsg:  "query user's videoList failed!"},
		})
		return
	}

	c.JSON(http.StatusOK, types.VideoListResponse{
		Response:  ErrNo.SuccessResp,
		VideoList: videoList,
	})
	log.Printf("获得用户发布过视频列表成功。uid:%d\n", videoListRequest.UserId)

}
