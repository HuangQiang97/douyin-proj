package controller

import (
	"douyin-proj/src/global/ErrNo"
	"douyin-proj/src/global/util"
	"douyin-proj/src/service"
	"douyin-proj/src/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RelationAction(c *gin.Context) {
	var relationRequest = types.RelationRequest{}
	if err := c.ShouldBind(&relationRequest); err != nil {
		c.JSON(http.StatusOK, types.RelationResponse{
			Response: types.Response{StatusCode: ErrNo.ParamInvalid, StatusMsg: err.Error()},
		})
		return
	}
	// 校验jwt token
	uId, err := util.VerifyToken(relationRequest.Token)
	if err != nil {
		c.JSON(http.StatusOK, types.UserResponse{
			Response: ErrNo.AuthFailedResp,
		})
		return
	}

	switch relationRequest.ActionType {
	case 1:
		if err := service.CreateRelation(uId, uint(relationRequest.ToUserId)); err != nil {
			c.JSON(http.StatusOK, types.RelationResponse{
				Response: types.Response{StatusCode: ErrNo.UnknownError, StatusMsg: err.Error()},
			})
			return
		}
	case 2:
		if err := service.DeleteRelation(uId, uint(relationRequest.ToUserId)); err != nil {
			c.JSON(http.StatusOK, types.RelationResponse{
				Response: types.Response{StatusCode: ErrNo.UnknownError, StatusMsg: err.Error()},
			})
			return
		}
	}

	c.JSON(http.StatusOK, types.RelationResponse{
		Response: ErrNo.SuccessResp,
	})
	return
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	var userFollowListRequest = types.UserFollowListRequest{}
	if err := c.ShouldBind(&userFollowListRequest); err != nil {
		c.JSON(http.StatusOK, types.UserFollowListResponse{
			Response: types.Response{StatusCode: ErrNo.ParamInvalid, StatusMsg: err.Error()},
		})
		return
	}
	// 校验jwt token
	uId, err := util.VerifyToken(userFollowListRequest.Token)
	if err != nil {
		c.JSON(http.StatusOK, types.UserResponse{
			Response: ErrNo.AuthFailedResp,
		})
		return
	}
	followList, err := service.GetFollowList(uint(userFollowListRequest.UserId), uId)
	if err != nil {
		c.JSON(http.StatusOK, types.UserFollowListResponse{
			Response: types.Response{StatusCode: ErrNo.UnknownError, StatusMsg: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, types.UserFollowListResponse{
		Response: ErrNo.SuccessResp,
		UserList: followList,
	})
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	var userFansListRequest = types.UserFansListRequest{}
	if err := c.ShouldBind(&userFansListRequest); err != nil {
		c.JSON(http.StatusOK, types.UserFansListResponse{
			Response: types.Response{StatusCode: ErrNo.ParamInvalid, StatusMsg: err.Error()},
		})
		return
	}

	// 校验jwt token
	uId, err := util.VerifyToken(userFansListRequest.Token)
	if err != nil {
		c.JSON(http.StatusOK, types.UserResponse{
			Response: ErrNo.AuthFailedResp,
		})
		return
	}

	fans, err := service.GetFansList(uint(userFansListRequest.UserId), uId)
	if err != nil {
		c.JSON(http.StatusOK, types.UserFansListResponse{
			Response: types.Response{StatusCode: ErrNo.UnknownError, StatusMsg: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, types.UserFansListResponse{
		Response: ErrNo.SuccessResp,
		UserList: fans,
	})
}
