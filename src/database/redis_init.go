package database

import "github.com/go-redis/redis"

var RedisClient *redis.Client

func InitRedis() {
	// 根据redis配置初始化一个客户端
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // redis地址
		Password: "",               // redis密码，没有则留空
		DB:       0,                // 默认数据库，默认是0
	})
}
