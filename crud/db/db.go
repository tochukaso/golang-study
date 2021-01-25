package db

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"omori.jp/env"
)

func GetDB() *gorm.DB {
	dsn := "gostudy:gostudy@/gostudy?parseTime=true"

	logFilePath := env.GetEnv().LogFilePath
	f, _ := os.Create(logFilePath)

	newLogger := logger.New(
		log.New(f, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,      // Slow SQL threshold
			LogLevel:      getSQLLogLevel(), // Log level
			Colorful:      false,            // Disable color
		},
	)

	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	log.SetOutput(f)
	return db
}

func getSQLLogLevel() logger.LogLevel {
	sqlLogLevel := env.GetEnv().SQLLogLevel

	var logLevel = logger.Silent
	switch sqlLogLevel {
	case "Silent":
		logLevel = logger.Silent
	case "Error":
		logLevel = logger.Error
	case "Warn":
		logLevel = logger.Warn
	case "Info":
		logLevel = logger.Info
	}
	return logLevel
}
