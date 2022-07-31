package database

// Init 初始数据库连接
func Init() error {
	if err := initMySQL(); err != nil {
		return err
	}

	return nil
}

// Close 关闭数据库连接
func Close() {
	db, _ := MySQLDb.DB()
	db.Close()
}
