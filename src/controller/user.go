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
	var userRegisterRequest = types.UserLoginRequest{}
	if err := c.ShouldBind(&userRegisterRequest); err != nil {
		c.JSON(http.StatusOK, types.UserLoginResponse{
			Response: ErrNo.ParamInvalidResp,
			Token:    "",
		})
		return
	}
	//username := c.Query("username")
	//password := c.Query("password")
	//if len(username) == 0 || len(password) == 0 {
	//	c.JSON(http.StatusOK, types.UserLoginResponse{
	//		Response: types.Response{StatusCode: 2, StatusMsg: "param error"},
	//		Token:    "",
	//	})
	//	return
	//}
	// user is existed
	userExisted, err := repository.GetUserByName(userRegisterRequest.UserName)
	if userExisted != nil {
		c.JSON(http.StatusOK, types.UserLoginResponse{
			Response:ErrNo.UserHasExistedResp,
			Token: "",
		})
		return
	}

	id, err := service.CreateUser(userRegisterRequest.UserName, userRegisterRequest.Password)
	if err != nil {
		c.JSON(http.StatusOK, types.UserLoginResponse{
			Response: types.Response{StatusCode: 3, StatusMsg: "create user failed"},
			Token:    "",
		})
		return
	}
	token, err := util.ReleaseToken(id)
	if err != nil {
		c.JSON(http.StatusOK, types.UserLoginResponse{
			Response: types.Response{StatusCode: 3, StatusMsg: "create user failed"},
			Token:    "",
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
			Response: ErrNo.ParamInvalidResp,
			Token:    "",
		})
		return
	}

	//username := c.Query("username")
	//password := c.Query("password")
	id, err := service.CheckUser(userLoginRequest.UserName, userLoginRequest.Password)
	if err != nil {
		c.JSON(http.StatusOK, types.UserLoginResponse{
			Response: types.Response{StatusCode: 4, StatusMsg: err.Error()},
			UserId:   0,
			Token:    "",
		})
		return
	}
	token, err := util.ReleaseToken(id)
	if err != nil {
		c.JSON(http.StatusOK, types.UserLoginResponse{
			Response: types.Response{StatusCode: 4, StatusMsg: "token init failed,login failed"},
			UserId:   0,
			Token:    "",
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
	var userInfoRequest =  types.UserInfoRequest{}
	if err := c.ShouldBind(&userInfoRequest); err != nil {
		c.JSON(http.StatusOK, types.FavoriteListResponse{
			Response: ErrNo.ParamInvalidResp,
		})
		return
	}
	uId, err := util.VerifyToken(userInfoRequest.Token)
	if err != nil{
		c.JSON(http.StatusOK, types.UserResponse{
			Response: ErrNo.AuthFailedResp,
			User: types.User{},
		})
		return
	}
	user, err := repository.GetUserById(uId)
	if err != nil{
		c.JSON(http.StatusOK, types.UserResponse{
			Response: types.Response{StatusCode: 4, StatusMsg: err.Error()},
			User: types.User{},
		})
	}
	if user == nil{
		c.JSON(http.StatusOK, types.UserResponse{
			Response: types.Response{StatusCode: 4, StatusMsg: "user id is unvalid"},
			User: types.User{} ,
		})
		return
	}
	c.JSON(http.StatusOK, types.UserResponse{
		Response: ErrNo.SuccessResp,
		User: types.User{
			Id: user.ID,
			Name: user.UserName,
			FollowCount: user.FollowCount,
			FollowerCount: user.FansCount,
			IsFollow: false, // todo: follow 关系查询
		},
	})




}
