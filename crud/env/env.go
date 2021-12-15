package env

import (
	"log"
	"strconv"

	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	DSN              string
	ProductImagePath string
	SmtpFrom         string
	SmtpHost         string
	SmtpPort         string
	SmtpPassword     string
	LogFilePath      string
	SQLLogLevel      string
	BaseURL          string
	CookieSSL        string
	CookieSameSite   string
}

func GetEnv() Env {
	var s Env
	err := envconfig.Process("golangstudy", &s)
	if err != nil {
		log.Fatal(err.Error())
	}
	return s
}

func GetCookieSSL() bool {
	b, err := strconv.ParseBool(GetEnv().CookieSSL)
	if err != nil {
		return false
	}
	return b
}
