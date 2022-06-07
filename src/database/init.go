package database

func Init() error {
	if err := initMySQL(); err != nil {
		return err
	}
	InitRedis()
	return nil
}

func Close() {
	db, _ := MySQLDb.DB()
	db.Close()
}
