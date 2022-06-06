package middleware

import (
	"douyin-proj/src/global/ErrNo"
	"douyin-proj/src/global/util"
	"douyin-proj/src/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func TokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, ok := c.GetQuery("token")
		if !ok {
			c.JSON(http.StatusOK, ErrNo.ParamInvalid)
			c.Abort()
			return
		}
		_, err := util.VerifyToken(token)
		if err != nil {
			c.JSON(http.StatusOK, types.Response{StatusCode: ErrNo.AuthFailed, StatusMsg: err.Error()})
			c.Abort()
			return
		}
		c.Next()
	}
}
