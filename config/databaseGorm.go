package config

import (
	"fmt"
	"os"

	mysql "gorm.io/driver/mysql"
	gorm "gorm.io/gorm"
)

var (
	db *gorm.DB
)

func GetDB() *gorm.DB {
	return db
}

//ConnectGorm Database Connection to Gorm V2
func ConnectGorm() {
	databaseConfig := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true&parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	)

	var err error
	db, err = gorm.Open(mysql.Open(databaseConfig), initConfig())

	if err != nil {
		fmt.Println(err)
		panic("Fail To Connect Database")
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(20)
}
