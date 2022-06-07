package main

import (
	"douyin-proj/src/config"
	"douyin-proj/src/database"
	"douyin-proj/src/repository"
	"douyin-proj/src/server"
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
	if err := config.Init(config.DefaultPath); err != nil {
		return err
	}
	if err := database.Init(); err != nil {
		return err
	}
	repository.Init()
	return nil
}
