package cache

import (
	"go_ws/config"
	"log"
	"github.com/go-redis/redis"
)

const (
       DB = config.REDIS_DB
       USERNAME = config.REDIS_USERNAME
       PASSWORD = config.REDIS_PASSWORD
       HOST = config.REDIS_HOST
       PORT = config.REDIS_PORT
)

var Client = new(redis.Client)

func init() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("error: %v", err)
		}
	}()
	Client = redis.NewClient(&redis.Options{
		Addr: HOST + ":" + PORT,
		Password: PASSWORD,
		DB: DB,
	})
}

func ClientClose()  {
	Client.Close()
}