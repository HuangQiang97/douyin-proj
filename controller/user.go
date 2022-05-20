package controller

import (
	"github.com/HuangQiang97/douyin-proj/entity"
	"github.com/HuangQiang97/douyin-proj/pkg/errno"
	"github.com/HuangQiang97/douyin-proj/pkg/util"
	"github.com/HuangQiang97/douyin-proj/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User entity.User `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	if len(username) == 0 || len(password) == 0{
		c.JSON(http.StatusOK, UserLoginResponse{
			Response:Response{StatusCode:errno.ParamErrCode,StatusMsg: errno.ParamErr.ErrMsg}})
		return
	}
	err := service.CreateUser(username, password)
	if err != nil{
		Err := errno.ConvertErr(err)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response:Response{StatusCode:Err.ErrCode,StatusMsg: Err.ErrMsg}})
		return
	}

	user, _ := service.GetUserByName(username)
	token, err := util.ReleaseToken(user.Id)
	if err != nil{
		Err := errno.ConvertErr(err)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response:Response{StatusCode:Err.ErrCode,StatusMsg: Err.ErrMsg}})
		return
	}
	c.JSON(http.StatusOK, UserLoginResponse{
		Response:Response{StatusCode: errno.SuccessCode, StatusMsg: errno.Success.ErrMsg},
		UserId: user.Id,
		Token: token,
	})
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	uId, err := service.CheckUser(username, password)
	if err != nil{
		Err := errno.ConvertErr(err)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response:Response{StatusCode:Err.ErrCode,StatusMsg: Err.ErrMsg}})
		return
	}
	token, err := util.ReleaseToken(uId)
	if err != nil{
		Err := errno.ConvertErr(err)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response:Response{StatusCode:Err.ErrCode,StatusMsg: Err.ErrMsg}})
		return
	}
	c.JSON(http.StatusOK, UserLoginResponse{
		Response:Response{StatusCode: errno.SuccessCode, StatusMsg: errno.Success.ErrMsg},
		UserId: uId,
		Token: token,
	})
}

func UserInfo(c *gin.Context) {
	token := c.Query("token")
	userIdStr := c.Query("user_id")
	if len(token) == 0 || len(userIdStr) == 0 {
		c.JSON(http.StatusOK, UserResponse{
			Response:Response{
				StatusCode: errno.ParamErrCode,
				StatusMsg: errno.ParamErr.ErrMsg,
			},
		})
		return
	}
	ownId , err := util.VerifyToken(token)
	if err != nil{
		c.JSON(http.StatusOK,UserResponse{
			Response:Response{
				StatusCode: errno.TokenInvalidErrCode,
				StatusMsg: errno.TokenInvalidErr.ErrMsg,
			},
		})
		return
	}

	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	user, err := service.GetUserAllInfo(userId)
	if err != nil{
		Err := errno.ConvertErr(err)
		c.JSON(http.StatusOK, UserResponse{
			Response:Response{StatusCode: Err.ErrCode, StatusMsg: Err.ErrMsg},
		})
		return
	}
	isFollow, err := service.CheckRelation(ownId, userId)
	if err != nil{
		Err := errno.ConvertErr(err)
		c.JSON(http.StatusOK, UserResponse{
			Response:Response{StatusCode: Err.ErrCode, StatusMsg: Err.ErrMsg},
		})
		return
	}
	user.IsFollow = isFollow
	c.JSON(http.StatusOK, UserResponse{
		Response:Response{StatusCode: errno.SuccessCode, StatusMsg: errno.Success.ErrMsg},
		User: *user,
	})
}
