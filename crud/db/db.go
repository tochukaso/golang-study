package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tochukaso/golang-study/env"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetDB() *gorm.DB {

	logFilePath := env.GetEnv().LogFilePath
	f, _ := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	newLogger := logger.New(
		log.New(f, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,      // Slow SQL threshold
			LogLevel:      getSQLLogLevel(), // Log level
			Colorful:      false,            // Disable color
		},
	)

	db, err := gorm.Open(mysql.Open(env.GetEnv().DSN), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		fmt.Println(err)
	}

	log.SetOutput(f)
	return db
}

func getSQLLogLevel() logger.LogLevel {
	sqlLogLevel := env.GetEnv().SQLLogLevel

	var logLevel = logger.Info
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
