package config

import (
	"gopkg.in/ini.v1"
	"log"
	"os"
	"time"
)

const DefaultPath string = "resources/application.ini"

//MySQLInfo 数据库连接信息
type MySQLInfo struct {
	TYPE         string
	USER         string
	PASSWORD     string
	DB_HOST      string
	DB_PORT      string
	DB_NAME      string
	CHARSET      string
	ParseTime    string
	MaxIdleConns int
	MaxOpenConns int
	Loc          string
}

// RedisInfo redis 连接信息
type RedisInfo struct {
	Addr     string
	Port     string
	Password string
	DB       string
}

// ServerInfo 对外暴露信息
type ServerInfo struct {
	HTTP_PORT string
	HTTP_HOST string
	MODE      string
}

var (
	MySQLConfig  *MySQLInfo
	ServerConfig *ServerInfo
	SecretKey    string
	Salt         string
	RedisConfig  *RedisInfo
)

// UrlPrefix 存储的图片和视频的链接
const PlayUrlPrefix = "http://10.192.58.230:8532/upload/video/"
const CoverUrlPrefix = "http://10.192.58.230:8532/upload/cover/"

// SavePrefix 视频图片保存路径
const VideoSavePrefix = "./upload/video/"
const CoverSavePrefix = "./upload/cover/"
const MaxMsgCount = 100

// RedisExpireTime redis缓存过期时间
const RedisExpireTime = 4 * time.Hour

// Init 初始化配置信息
func Init(path string) error {
	cfg, err := ini.Load(path)
	if err != nil {
		return err
	}
	if err := initMysqlConfig(cfg); err != nil {
		return err
	}
	if err := initServerConfig(cfg); err != nil {
		return err
	}
	initLog()
	SecretKey = cfg.Section("jwt").Key("secretKey").String()
	Salt = cfg.Section("crypto").Key("salt").String()
	initRedis(cfg)
	log.Println("初始化配置成功")
	return nil
}

// initServerConfig 初始化服务器暴露信息
func initServerConfig(cfg *ini.File) error {
	ServerConfig = new(ServerInfo)
	return cfg.Section("server").MapTo(ServerConfig)
}

// initMysqlConfig 初始化数据库配置
func initMysqlConfig(cfg *ini.File) error {
	MySQLConfig = new(MySQLInfo)
	return cfg.Section("mysql").MapTo(MySQLConfig)
}

// initLog 初始化日志系统
func initLog() {
	file := "./log.txt"
	logFile, _ := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	log.SetOutput(logFile) // 将文件设置为log输出的文件
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
}

// initRedis 初始化redis信息
func initRedis(cfg *ini.File) error {
	RedisConfig = new(RedisInfo)
	return cfg.Section("redis").MapTo(RedisConfig)
}
