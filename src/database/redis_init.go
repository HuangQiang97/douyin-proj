package database

import (
	"douyin-proj/src/config"
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
)

var RedisClient *redis.Client

func InitRedis() {
	// 根据redis配置初始化一个客户端
	db, _ := strconv.Atoi(config.RedisConfig.DB)
	addr := fmt.Sprintf("%s:%s", config.RedisConfig.Addr, config.RedisConfig.Port)
	password := config.RedisConfig.Password
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,     // redis地址
		Password: password, // redis密码，没有则留空
		DB:       db,       // 默认数据库，默认是0
	})
}
