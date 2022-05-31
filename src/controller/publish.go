package controller

import (
	"douyin-proj/src/global/ErrNo"
	"douyin-proj/src/types"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func Publish(c *gin.Context) {
	var publishRequest = types.PublishRequest{}
	if err := c.ShouldBind(&publishRequest); err != nil {
		c.JSON(http.StatusOK, types.PublishResponse{StatusCode: ErrNo.ParamInvalid, StatusMsg: err.Error()})
		return
	}
	//destfile, err := os.OpenFile("./upload/"+publishRequest.Data.Filename, os.O_CREATE, 0666)
	// 如果已经存在，会将文件清空
	// TODO：需要设计文件存储的路径，以及文件名
	destfile, err := os.Create("./upload/" + publishRequest.Data.Filename)
	if err != nil {
		// 创建文件失败，原因可能是1.路径不存在2.权限不足3.打开文件数量超过上限4.磁盘空间不足
		fmt.Println(err)
		c.JSON(http.StatusOK, ErrNo.ParamInvalid)
		return
	}
	// 将request中的文件保存到目标文件
	fmt.Println(destfile.Name())
	err = c.SaveUploadedFile(publishRequest.Data, destfile.Name())
	if err != nil {
		fmt.Println(err)
	}

	c.JSON(http.StatusOK, ErrNo.SuccessResp)
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {

}
