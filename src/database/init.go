package database

func Init() error {
	if err := initMySQL(); err != nil {
		return err
	}
	return nil
}

func Close() {
	db, _ := MySQLDb.DB()
	db.Close()
}
