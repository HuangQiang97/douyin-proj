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

// RelationAction 添加/删除关注
func RelationAction(c *gin.Context) {
	var relationRequest = types.RelationRequest{}
	if err := c.ShouldBind(&relationRequest); err != nil {
		log.Printf("反序列化增删关注请求失败。err:%s\n", err)
		c.JSON(http.StatusOK, types.RelationResponse{
			Response: types.Response{StatusCode: ErrNo.ParamInvalid, StatusMsg: err.Error()},
		})
		return
	}
	// 校验jwt token
	uId, err := util.VerifyToken(relationRequest.Token)
	if err != nil {
		log.Printf("登录失败，err:%s\n", err)
		c.JSON(http.StatusOK, types.UserResponse{
			Response: ErrNo.AuthFailedResp,
		})
		return
	}

	switch relationRequest.ActionType {
	// 添加关注
	case 1:
		if err := service.CreateRelation(uId, uint(relationRequest.ToUserId)); err != nil {
			c.JSON(http.StatusOK, types.RelationResponse{
				Response: types.Response{StatusCode: ErrNo.RelationAddFailed, StatusMsg: err.Error()},
			})
			return
		}
		log.Printf("添加关注成功，uid:%d,follwId:%d\n", uId, relationRequest.ToUserId)
	// 删除关注
	case 2:
		if err := service.DeleteRelation(uId, uint(relationRequest.ToUserId)); err != nil {
			c.JSON(http.StatusOK, types.RelationResponse{
				Response: types.Response{StatusCode: ErrNo.RelationDeleteFailed, StatusMsg: err.Error()},
			})
			return
		}
		log.Printf("删除关注成功，uid:%d,follwId:%d\n", uId, relationRequest.ToUserId)
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
		log.Printf("反序列化获取关注列表请求失败。err:%s\n", err)
		c.JSON(http.StatusOK, types.UserFollowListResponse{
			Response: types.Response{StatusCode: ErrNo.ParamInvalid, StatusMsg: err.Error()},
		})
		return
	}
	// 校验jwt token
	uId, err := util.VerifyToken(userFollowListRequest.Token)
	if err != nil {
		log.Printf("登录失败，err:%s\n", err)
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
	log.Printf("获取关注列表成功，uid:%d\n", uId)

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
