package controller

import (
	"douyin-proj/src/config"
	"douyin-proj/src/server/middleware"
	"douyin-proj/src/service"
	"douyin-proj/src/types"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

// Register 注册
func Register(c *gin.Context) {
	// 参数绑定与校验
	var userRegisterRequest = types.UserLoginRequest{}
	if err := c.ShouldBind(&userRegisterRequest); err != nil {
		log.Printf("反序列化注册请求失败。err:%s\n", err)
		c.JSON(http.StatusOK, types.UserLoginResponse{
			Response: types.Response{StatusCode: config.ParamInvalid, StatusMsg: err.Error()},
		})
		return
	}
	// 空白内容判断
	if len(strings.TrimSpace(userRegisterRequest.UserName)) == 0 || len(strings.TrimSpace(userRegisterRequest.Password)) == 0 {
		log.Printf("密码或用户名为空白字符串。\n")
		c.JSON(http.StatusOK, types.UserLoginResponse{
			Response: types.Response{StatusCode: config.ParamInvalid, StatusMsg: "密码或用户名非法"},
		})
		return
	}
	// 创建用户
	id, err := service.CreateUser(userRegisterRequest.UserName, userRegisterRequest.Password)
	if err != nil {
		c.JSON(http.StatusOK, types.UserLoginResponse{
			Response: config.UserHasExistedResp,
		})
		return
	}
	// 生成jwt token
	token, err := middleware.ReleaseToken(id)
	if err != nil {
		c.JSON(http.StatusOK, types.UserLoginResponse{
			Response: types.Response{
				StatusCode: config.UnknownError,
				StatusMsg:  "Token init failed,register success.Please try to login!"},
		})
		return
	}
	c.JSON(http.StatusOK, types.UserLoginResponse{
		Response: config.SuccessResp,
		UserId:   int64(id),
		Token:    token,
	})
	log.Printf("用户注册成功，uid:%d\n", id)
}

// Login 登录
func Login(c *gin.Context) {
	var userLoginRequest = types.UserLoginRequest{}
	if err := c.ShouldBind(&userLoginRequest); err != nil {
		log.Printf("反序列化登录请求失败。err:%s\n", err)
		c.JSON(http.StatusOK, types.UserLoginResponse{
			Response: types.Response{StatusCode: config.ParamInvalid, StatusMsg: err.Error()},
		})
		return
	}
	// 检查用户名和密码
	id, err := service.CheckUser(userLoginRequest.UserName, userLoginRequest.Password)
	if err != nil {
		errCode := int32(id)
		if errCode == config.UserNotExisted {
			log.Printf("用户不存在。username:%s\n", userLoginRequest.UserName)
			c.JSON(http.StatusOK, types.UserLoginResponse{
				Response: config.UserNotExistedResp,
			})
		}
		if errCode == config.WrongPassword {
			log.Printf("密码错误。username:%s\n", userLoginRequest.UserName)
			c.JSON(http.StatusOK, types.UserLoginResponse{
				Response: config.WrongPasswordResp,
			})
		}
		return
	}
	// 发放token
	token, err := middleware.ReleaseToken(id)
	if err != nil {
		c.JSON(http.StatusOK, types.UserLoginResponse{
			Response: types.Response{StatusCode: config.UnknownError, StatusMsg: "Token init failed,login failed!"},
		})
		return
	}
	c.JSON(http.StatusOK, types.UserLoginResponse{
		Response: config.SuccessResp,
		UserId:   int64(id),
		Token:    token,
	})
	log.Printf("用户登录成功。uid:%d\n", id)
}

// UserInfo 获取用户信息
func UserInfo(c *gin.Context) {
	var userInfoRequest = types.UserInfoRequest{}
	if err := c.ShouldBind(&userInfoRequest); err != nil {
		c.JSON(http.StatusOK, types.FavoriteListResponse{
			Response: types.Response{StatusCode: config.ParamInvalid, StatusMsg: err.Error()},
		})
		return
	}
	// 校验jwt token
	uId, err := middleware.VerifyToken(userInfoRequest.Token)
	if err != nil {
		c.JSON(http.StatusOK, types.UserResponse{
			Response: config.AuthFailedResp,
		})
		return
	}
	// 查询用户
	user, err := service.GetUserInfo(userInfoRequest.UserId, uId)
	if err != nil {
		c.JSON(http.StatusOK, types.UserResponse{
			Response: config.UserNotExistedResp,
		})
		return
	}
	c.JSON(http.StatusOK, types.UserResponse{
		Response: config.SuccessResp,
		User: types.User{
			Id:            user.Id,
			Name:          user.Name,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FollowerCount,
			IsFollow:      user.IsFollow,
		},
	})
}
