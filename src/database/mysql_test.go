package database

import (
	"douyin-proj/src/config"
	"testing"
)

func TestMain(m *testing.M) {
	if err := config.Init("../../resources/application.ini"); err != nil {
		panic(m)
	}
	m.Run()
}

func Test_initMySQL(t *testing.T) {
	if err := initMySQL(); err != nil {
		t.Errorf("initMySQL() error = %v", err)
	}
	db, _ := MySQLDb.DB()
	defer db.Close()
	if err := db.Ping(); err != nil {
		t.Errorf("ping error = %v", err)
	}
}
