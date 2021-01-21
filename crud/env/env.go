package env

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	ProductImagePath string
	SmtpFrom         string
	SmtpHost         string
	SmtpPort         string
	SmtpPassword     string
}

func GetEnv() Env {
	var s Env
	err := envconfig.Process("golangstudy", &s)
	if err != nil {
		log.Fatal(err.Error())
	}
	return s
}
