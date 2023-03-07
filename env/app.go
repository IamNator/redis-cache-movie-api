package env

import (
	"os"
	"sync"
)

type (
	appEnv struct {
		PORT         string `json:"port"`
		REDIS_URL    string `json:"redis_url"`
		POSTGRES_URL string `json:"postgres_url"`
		HOST_MACHINE string `json:"host_machine"`
	}
)

var (
	sO         sync.Once
	defaultEnv appEnv
)

func Init() error {
	sO.Do(func() {
		defaultEnv.PORT = os.Getenv("PORT")
		defaultEnv.REDIS_URL = os.Getenv("REDISCLOUD_URL")
		defaultEnv.POSTGRES_URL = os.Getenv("DATABASE_URL")
		defaultEnv.HOST_MACHINE = os.Getenv("HOST_MACHINE")
	})

	return nil
}

func Get() appEnv {
	return defaultEnv
}
