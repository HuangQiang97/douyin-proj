package config

import (
	"gopkg.in/ini.v1"
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

type ServerInfo struct {
	HTTP_PORT string
	HTTP_HOST string
	MODE      string
}

var (
	MySQLConfig  *MySQLInfo
	ServerConfig *ServerInfo
	SecretKey    string
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
	SecretKey = cfg.Section("jwt").Key("secretKey").String()
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
