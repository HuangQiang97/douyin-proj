package database

import (
	"douyin-proj/src/config"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

var MySQLDb *gorm.DB

func initMySQL() error {
	//dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		config.MySQLConfig.USER,
		config.MySQLConfig.PASSWORD,
		config.MySQLConfig.DB_HOST,
		config.MySQLConfig.DB_PORT,
		config.MySQLConfig.DB_NAME,
		config.MySQLConfig.CHARSET,
		config.MySQLConfig.ParseTime,
		config.MySQLConfig.Loc,
	)
	var err error
	MySQLDb, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			NoLowerCase:   false,
		},
		PrepareStmt: true,
	})
	if err != nil {
		return err
	}
	if config.ServerConfig.MODE == "debug" {
		MySQLDb.Logger.LogMode(logger.Silent)
	}
	db, _ := MySQLDb.DB()
	db.SetMaxIdleConns(config.MySQLConfig.MaxIdleConns)
	db.SetMaxOpenConns(config.MySQLConfig.MaxOpenConns)
	db.SetConnMaxIdleTime(30 * time.Minute)

	return nil
}
