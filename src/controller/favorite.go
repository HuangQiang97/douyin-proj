package controller

import (
	"douyin-proj/src/global/ErrNo"
	"douyin-proj/src/global/util"
	"douyin-proj/src/service"
	"douyin-proj/src/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

// FavoriteAction mark a video as favorite (action_type == 1) or undo it (action_type == 2)
func FavoriteAction(c *gin.Context) {
	// parse parameters
	var favoriteRequest = types.FavoriteRequest{}
	if err := c.ShouldBind(&favoriteRequest); err != nil {
		c.JSON(http.StatusOK, types.FavoriteResponse{
			Response: types.Response{StatusCode: ErrNo.ParamInvalid, StatusMsg: err.Error()},
		})
		return
	}
	// check token
	uId, err := util.VerifyToken(favoriteRequest.Token)
	if err != nil {
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
	case 2: // undo favorite
		if err := service.UndoFavorite(uId, favoriteRequest.VideoId); err != nil {
			c.JSON(http.StatusOK, types.FavoriteResponse{
				Response: ErrNo.NotInFavoriteResp,
			})
			return
		}
	}
	c.JSON(http.StatusOK, types.FavoriteResponse{
		Response: ErrNo.SuccessResp,
	})
	return
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	var favoriteListRequest = types.FavoriteListRequest{}
	if err := c.ShouldBind(&favoriteListRequest); err != nil {
		c.JSON(http.StatusOK, types.FavoriteListResponse{
			Response: types.Response{StatusCode: ErrNo.ParamInvalid, StatusMsg: err.Error()},
		})
		return
	}

	// check token
	uId, err := util.VerifyToken(favoriteListRequest.Token)
	if err != nil {
		c.JSON(http.StatusOK, types.FavoriteListResponse{
			Response: ErrNo.AuthFailedResp,
		})
		return
	}

	videoList, err := service.GetFavoriteVideoListByUserId(favoriteListRequest.UserId, uId)
	if err != nil {
		c.JSON(http.StatusOK, types.FavoriteListResponse{
			Response: types.Response{StatusCode: ErrNo.UnknownError, StatusMsg: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, types.FavoriteListResponse{
		Response:  ErrNo.SuccessResp,
		VideoList: videoList,
	})
}
