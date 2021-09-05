package config

import (
	"log"
	"os"
	"time"

	gorm "gorm.io/gorm"
	logger "gorm.io/gorm/logger"
	schema "gorm.io/gorm/schema"
)

//Generated By github.com/david-yappeter/GormCrudGenerator

//initConfig Initialize Config
func initConfig() *gorm.Config {
	return &gorm.Config{
		Logger:         initLog(),
		NamingStrategy: initNamingStrategy(),
	}
}

//initLog Connection Log Configuration
func initLog() logger.Interface {
	newLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		Colorful:      true,
		LogLevel:      logger.Info,
		SlowThreshold: time.Second,
	})
	return newLogger
}

//initNamingStrategy Init NamingStrategy
func initNamingStrategy() *schema.NamingStrategy {
	return &schema.NamingStrategy{
		SingularTable: true,
		TablePrefix:   "",
	}
}
