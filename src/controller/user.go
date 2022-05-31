package controller

import (
	"douyin-proj/src/global/ErrNo"
	"douyin-proj/src/global/util"
	"douyin-proj/src/repository"
	"douyin-proj/src/service"
	"douyin-proj/src/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Register(c *gin.Context) {
	// 参数绑定与校验
	var userRegisterRequest = types.UserLoginRequest{}
	if err := c.ShouldBind(&userRegisterRequest); err != nil {
		c.JSON(http.StatusOK, types.UserLoginResponse{
			Response: types.Response{StatusCode: ErrNo.ParamInvalid, StatusMsg: err.Error()},
		})
		return
	}
	// 创建用户
	id, err := service.CreateUser(userRegisterRequest.UserName, userRegisterRequest.Password)
	if err != nil {
		c.JSON(http.StatusOK, types.UserLoginResponse{
			Response: ErrNo.UserHasExistedResp,
		})
		return
	}
	// 生成jwt token
	token, err := util.ReleaseToken(id)
	if err != nil {
		c.JSON(http.StatusOK, types.UserLoginResponse{
			Response: types.Response{
				StatusCode: ErrNo.UnknownError,
				StatusMsg:  "Token init failed,register success.Please try to login!"},
		})
		return
	}

	c.JSON(http.StatusOK, types.UserLoginResponse{
		Response: ErrNo.SuccessResp,
		UserId:   int64(id),
		Token:    token,
	})
}

func Login(c *gin.Context) {
	var userLoginRequest = types.UserLoginRequest{}
	if err := c.ShouldBind(&userLoginRequest); err != nil {
		c.JSON(http.StatusOK, types.UserLoginResponse{
			Response: types.Response{StatusCode: ErrNo.ParamInvalid, StatusMsg: err.Error()},
		})
		return
	}
	// 检查用户名和密码
	id, err := service.CheckUser(userLoginRequest.UserName, userLoginRequest.Password)
	if err != nil {
		errCode := int32(id)
		if errCode == ErrNo.UserNotExisted {
			c.JSON(http.StatusOK, types.UserLoginResponse{
				Response: ErrNo.UserNotExistedResp,
			})
		}
		if errCode == ErrNo.WrongPassword {
			c.JSON(http.StatusOK, types.UserLoginResponse{
				Response: ErrNo.WrongPasswordResp,
			})
		}
	}

	token, err := util.ReleaseToken(id)
	if err != nil {
		c.JSON(http.StatusOK, types.UserLoginResponse{
			Response: types.Response{StatusCode: ErrNo.UnknownError, StatusMsg: "Token init failed,login failed!"},
		})
		return
	}
	c.JSON(http.StatusOK, types.UserLoginResponse{
		Response: ErrNo.SuccessResp,
		UserId:   int64(id),
		Token:    token,
	})
}

func UserInfo(c *gin.Context) {
	var userInfoRequest = types.UserInfoRequest{}
	if err := c.ShouldBind(&userInfoRequest); err != nil {
		c.JSON(http.StatusOK, types.FavoriteListResponse{
			Response: ErrNo.ParamInvalidResp,
		})
		return
	}
	// 校验jwt token
	uId, err := util.VerifyToken(userInfoRequest.Token)
	if err != nil {
		c.JSON(http.StatusOK, types.UserResponse{
			Response: ErrNo.AuthFailedResp,
		})
		return
	}
	// 查询用户
	user, err := repository.GetUserById(uId)
	if err != nil {
		c.JSON(http.StatusOK, types.UserResponse{
			Response: ErrNo.UserNotExistedResp,
		})
	}

	c.JSON(http.StatusOK, types.UserResponse{
		Response: ErrNo.SuccessResp,
		User: types.User{
			Id:            user.ID,
			Name:          user.UserName,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FansCount,
			IsFollow:      false, // todo: follow 关系查询
		},
	})

}
