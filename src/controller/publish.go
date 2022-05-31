package controller

import (
	"douyin-proj/src/global/ErrNo"
	"douyin-proj/src/global/util"
	"douyin-proj/src/repository"
	"douyin-proj/src/service"
	"douyin-proj/src/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
)

func Publish(c *gin.Context) {
	var publishRequest = types.PublishRequest{}
	if err := c.ShouldBind(&publishRequest); err != nil {
		c.JSON(http.StatusOK, types.PublishResponse{StatusCode: ErrNo.ParamInvalid, StatusMsg: err.Error()})
		return
	}

	// 校验jwt token
	uId, err := util.VerifyToken(publishRequest.Token)
	if err != nil {
		c.JSON(http.StatusOK, types.UserResponse{
			Response: ErrNo.AuthFailedResp,
		})
		return
	}

	// 如果已经存在，会将文件清空
	// TODO：需要设计文件存储的路径，以及文件名
	destfile, err := os.Create("./upload/" + publishRequest.Data.Filename)
	if err != nil {
		// 创建文件失败，原因可能是1.路径不存在2.权限不足3.打开文件数量超过上限4.磁盘空间不足
		c.JSON(http.StatusOK, ErrNo.VideoUploadFailedResp)
		return
	}
	// 将request中的文件保存到目标文件
	err = c.SaveUploadedFile(publishRequest.Data, destfile.Name())
	if err != nil {
		c.JSON(http.StatusOK, ErrNo.VideoUploadFailedResp)
		return
	}
	// TODO : url
	video := repository.Video{
		AuthorID:      uId,
		PlayUrl:       "",
		CoverUrl:      "",
		FavoriteCount: 0,
		CommentCount:  0,
		Title:         publishRequest.Title,
		CreatedAt:     uint64(time.Now().Unix()),
	}
	if err = repository.CreateVideo(&video); err != nil {
		c.JSON(http.StatusOK, ErrNo.VideoUploadFailedResp)
		return
	}

	c.JSON(http.StatusOK, ErrNo.SuccessResp)
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	var videoListRequest = types.VideoListRequest{}
	if err := c.ShouldBind(&videoListRequest); err != nil {
		c.JSON(http.StatusOK, types.PublishResponse{StatusCode: ErrNo.ParamInvalid, StatusMsg: err.Error()})
		return
	}
	// 校验jwt token
	id, err := util.VerifyToken(videoListRequest.Token)
	if err != nil {
		c.JSON(http.StatusOK, types.VideoListResponse{
			Response: ErrNo.AuthFailedResp,
		})
		return
	}

	videoList, err := service.GetVideoList(id)
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
}
