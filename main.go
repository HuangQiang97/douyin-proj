package main

import (
	"douyin-proj/src/config"
	"douyin-proj/src/database"
	"douyin-proj/src/repository"
	"douyin-proj/src/server"
	"douyin-proj/src/server/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := Init(); err != nil {
		panic(err)
	}

	defer func() {
		database.Close()
	}()
	r := gin.Default()
	server.Run(r)

}

func Init() error {
	config.Init(config.DefaultPath)
	database.Init()
	repository.Init()
	go middleware.InitFfmpeg()
	middleware.InitRedis()
	return nil
}
