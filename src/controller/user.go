package controller

import (
	"douyin-proj/src/global/util"
	"douyin-proj/src/service"
	"douyin-proj/src/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	if len(username) == 0 || len(password) == 0 {
		c.JSON(http.StatusOK, types.UserRegisterResponse{
			Response: types.Response{StatusCode: 2, StatusMsg: "param error"},
			Token:    "",
		})
		return
	}
	id, err := service.CreateUser(username, password)
	if err != nil {
		c.JSON(http.StatusOK, types.UserRegisterResponse{
			Response: types.Response{StatusCode: 3, StatusMsg: "create user failed"},
			Token:    "",
		})
		return
	}
	token, err := util.ReleaseToken(id)
	if err != nil {
		c.JSON(http.StatusOK, types.UserRegisterResponse{
			Response: types.Response{StatusCode: 3, StatusMsg: "create user failed"},
			Token:    "",
		})
		return
	}

	c.JSON(http.StatusOK, types.UserRegisterResponse{
		Response: types.Response{StatusCode: 0, StatusMsg: "success"},
		UserId:   int64(id),
		Token:    token,
	})
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	id, err := service.CheckUser(username, password)
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
	c.JSON(http.StatusOK, types.UserRegisterResponse{
		Response: types.Response{StatusCode: 0, StatusMsg: "success"},
		UserId:   int64(id),
		Token:    token,
	})
}

func UserInfo(c *gin.Context) {

}
