package main

import (
	"github.com/HuangQiang97/douyin-proj/dal"
	"github.com/gin-gonic/gin"
)

func init(){
	dal.Init()
}

func main() {
	r := gin.Default()

	initRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
