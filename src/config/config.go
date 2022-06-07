package config

import (
	"gopkg.in/ini.v1"
	"log"
	"os"
)

const DefaultPath string = "resources/application.ini"

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
type RedisInfo struct {
	Addr     string
	Port     string
	Password string
	DB       string
}

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

func initServerConfig(cfg *ini.File) error {
	ServerConfig = new(ServerInfo)
	return cfg.Section("server").MapTo(ServerConfig)
}

func initMysqlConfig(cfg *ini.File) error {
	MySQLConfig = new(MySQLInfo)
	return cfg.Section("mysql").MapTo(MySQLConfig)
}
func initLog() {
	file := "./log.txt"
	logFile, _ := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	log.SetOutput(logFile) // 将文件设置为log输出的文件
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
}

func initRedis(cfg *ini.File) error {
	RedisConfig = new(RedisInfo)
	return cfg.Section("redis").MapTo(RedisConfig)
}
