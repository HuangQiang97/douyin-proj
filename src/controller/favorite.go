package controller

import (
	"douyin-proj/src/global/ErrNo"
	"douyin-proj/src/repository"
	"douyin-proj/src/types"
	"fmt"
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

	f := repository.Favorite{
		UserID:  favoriteRequest.UserId,
		VideoID: favoriteRequest.VideoId,
	}

	switch favoriteRequest.ActionType {
	case 1: // do favorite
		if err := repository.CreateFavorite(&f); err != nil {
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
		if err := repository.UndoFavorite(&f); err != nil {
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

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	var favoriteListRequest = types.FavoriteListRequest{}
	if err := c.ShouldBind(&favoriteListRequest); err != nil {
		c.JSON(http.StatusOK, types.FavoriteListResponse{
			Response: ErrNo.ParamInvalidResp,
		})
		return
	}

	// TODO check token

	videoIds, err := repository.GetFavoriteVideoIdsByUserId(favoriteListRequest.UserId)
	fmt.Println(videoIds)
	if err != nil {
		c.JSON(http.StatusOK, types.FavoriteListResponse{
			Response: ErrNo.UnknownErrorResp,
		})
		return
	} else {

		// return if no favorite video
		if len(videoIds) == 0 {
			c.JSON(http.StatusOK, types.FavoriteListResponse{
				Response:  ErrNo.SuccessResp,
				VideoList: []types.Video{},
			})
			return
		}

		// get repository.Video by video IDs
		// 注意：如果video表中不存在对应的video_id，那么会跳过此记录
		// 因此，可能出现video ID数量与此处返回的Video对象数量不一致的情况
		videoPtrList, err := repository.GetVideosByIds(videoIds)
		if err != nil {
			c.JSON(http.StatusOK, types.FavoriteListResponse{
				Response: ErrNo.UnknownErrorResp,
			})
			return
		}

		videoList := make([]types.Video, len(videoPtrList))
		for _, videoPtr := range videoPtrList {
			author, err := repository.GetUserById(videoPtr.AuthorID)
			if err != nil {
				c.JSON(http.StatusOK, types.FavoriteListResponse{
					Response: ErrNo.UnknownErrorResp,
				})
				return
			}

			// convert repository.Video to types.Video
			video := types.Video{
				Id: videoPtr.ID,
				Author: types.User{ // convert repository.User to types.User
					Id:            author.ID,
					Name:          author.UserName,
					FollowCount:   author.FollowCount,
					FollowerCount: author.FansCount,
					IsFollow:      false, // TODO: implement IsFollow checker
				},
				PlayUrl:       videoPtr.PlayUrl,
				CoverUrl:      videoPtr.CoverUrl,
				FavoriteCount: videoPtr.FavoriteCount,
				CommentCount:  videoPtr.CommentCount,
				IsFavorite: repository.IsFavorite(&repository.Favorite{
					UserID:  favoriteListRequest.UserId,
					VideoID: videoPtr.ID,
				}),
				Title: videoPtr.Title,
			}
			videoList = append(videoList, video)
		}
		c.JSON(http.StatusOK, types.FavoriteListResponse{
			Response:  ErrNo.SuccessResp,
			VideoList: videoList,
		})
		return
	}

}
