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

// FavoriteAction mark a video as favorite (action_type == 1) or undo it (action_type == 2)
func FavoriteAction(c *gin.Context) {
	// parse parameters
	var favoriteRequest = types.FavoriteRequest{}
	if err := c.ShouldBind(&favoriteRequest); err != nil {
		log.Printf("反序列化点赞操作请求失败。err:%s\n", err)
		c.JSON(http.StatusOK, types.FavoriteResponse{
			Response: types.Response{StatusCode: ErrNo.ParamInvalid, StatusMsg: err.Error()},
		})
		return
	}
	// check token
	uId, err := util.VerifyToken(favoriteRequest.Token)
	if err != nil {
		log.Printf("登录失败，err:%s\n", err)
		c.JSON(http.StatusOK, types.FavoriteResponse{
			Response: ErrNo.AuthFailedResp,
		})
		return
	}

	switch favoriteRequest.ActionType {
	case 1: // do favorite
		if err := service.AddFavorite(uId, favoriteRequest.VideoId); err != nil {
			c.JSON(http.StatusOK, types.FavoriteResponse{
				Response: ErrNo.DuplicateFavoriteResp,
			})
			return
		}
		log.Printf("添加点赞记录成功。 uid:%d,videoId:%d\n", uId, favoriteRequest.VideoId)
	case 2: // undo favorite
		if err := service.UndoFavorite(uId, favoriteRequest.VideoId); err != nil {
			c.JSON(http.StatusOK, types.FavoriteResponse{
				Response: ErrNo.NotInFavoriteResp,
			})
			return
		}
		log.Printf("删除点赞记录成功。 uid:%d,videoId:%d\n", uId, favoriteRequest.VideoId)
	}
	c.JSON(http.StatusOK, types.FavoriteResponse{
		Response: ErrNo.SuccessResp,
	})
	return
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	// parse parameters
	var favoriteListRequest = types.FavoriteListRequest{}
	if err := c.ShouldBind(&favoriteListRequest); err != nil {
		log.Printf("反序列化获取点赞过视频操作请求失败。err:%s\n", err)
		c.JSON(http.StatusOK, types.FavoriteListResponse{
			Response: types.Response{StatusCode: ErrNo.ParamInvalid, StatusMsg: err.Error()},
		})
		return
	}

	// check token
	uId, err := util.VerifyToken(favoriteListRequest.Token)
	if err != nil {
		log.Printf("登录失败，err:%s\n", err)
		c.JSON(http.StatusOK, types.FavoriteListResponse{
			Response: ErrNo.AuthFailedResp,
		})
		return
	}

	// 获得用户点赞过视频
	videoList, err := service.GetFavoriteVideoListByUserId(favoriteListRequest.UserId, uId)
	if err != nil {
		log.Printf("获得用户点赞过视频失败。 uid:%d,err:%s\n", favoriteListRequest.UserId, err)
		c.JSON(http.StatusOK, types.FavoriteListResponse{
			Response: types.Response{StatusCode: ErrNo.UnknownError, StatusMsg: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, types.FavoriteListResponse{
		Response:  ErrNo.SuccessResp,
		VideoList: videoList,
	})
	log.Printf("获得用户点赞过视频成功。 uid:%d\n", favoriteListRequest.UserId)
}
