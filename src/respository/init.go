package respository

import (
	"douyin-proj/src/database"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	DB = database.MySQLDb
}
