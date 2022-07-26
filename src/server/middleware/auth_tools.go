package middleware

import (
	"douyin-proj/src/config"
	"douyin-proj/src/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func TokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, ok := c.GetQuery("token")
		if !ok {
			c.JSON(http.StatusOK, config.ParamInvalid)
			c.Abort()
			return
		}
		_, err := VerifyToken(token)
		if err != nil {
			c.JSON(http.StatusOK, types.Response{StatusCode: config.AuthFailed, StatusMsg: err.Error()})
			c.Abort()
			return
		}
		c.Next()
	}
}
