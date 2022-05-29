package controller

import (
	"douyin-proj/src/global/ErrNo"
	"douyin-proj/src/respository"
	"douyin-proj/src/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

// FavoriteAction mark a video as favorite (action_type == 1) or undo it (action_type == 2)
func FavoriteAction(c *gin.Context) {
	var favoriteRequest = types.FavoriteRequest{}
	if err := c.ShouldBind(&favoriteRequest); err != nil {
		c.JSON(http.StatusOK, types.FavoriteResponse{
			Response: ErrNo.ParamInvalidResp,
		})
		return
	}

	// TODO: check token
	// if _, err := util.VerifyToken(favoriteRequest.Token); err != nil {
	// 	c.JSON(http.StatusOK, types.FavoriteResponse{
	// 		Response: ErrNo.NotSignedInResp,
	// 	})
	// 	return
	// }

	f := respository.Favorite{
		UserID:  favoriteRequest.UserId,
		VideoID: favoriteRequest.VideoId,
	}

	switch favoriteRequest.ActionType {
	case 1: // do favorite
		if err := respository.CreateFavorite(&f); err != nil {
			c.JSON(http.StatusOK, types.FavoriteResponse{
				Response: ErrNo.DuplicateFavoriteResp,
			})
			return
		} else {
			c.JSON(http.StatusOK, types.FavoriteResponse{
				Response: ErrNo.SuccessResp,
			})
			return
		}
	case 2: // undo favorite
		if err := respository.UndoFavorite(&f); err != nil {
			c.JSON(http.StatusOK, types.FavoriteResponse{
				Response: ErrNo.NotInFavoriteResp,
			})
			return
		} else {
			c.JSON(http.StatusOK, types.FavoriteResponse{
				Response: ErrNo.SuccessResp,
			})
			return
		}
	}
}

// FavoriteList all users have same favorite video list TODO
func FavoriteList(c *gin.Context) {

}
