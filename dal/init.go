package dal

import (
	"github.com/HuangQiang97/douyin-proj/pkg/constant"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	db, err := gorm.Open(mysql.Open(constant.MySQLDefaultDSN),&gorm.Config{
		PrepareStmt: true,
		SkipDefaultTransaction: true,
	})
	DB = db
	if err != nil{
		panic(err)
	}
	if !DB.Migrator().HasTable(constant.UserTableName){
		err = DB.Migrator().CreateTable(&User{})
		if err != nil{
			panic(err)
		}
	}
	if !DB.Migrator().HasTable(constant.RelationTableName){
		err = DB.Migrator().CreateTable(&Relation{})
		if err != nil{
			panic(err)
		}
	}
}

