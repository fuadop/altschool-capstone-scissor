package model

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func init() {
	url := os.Getenv("REDIS_URL")
	opt, err := redis.ParseURL(url)
	if err != nil {
		log.Fatalln(err)
	}

	rdb = redis.NewClient(opt)

	_, err = rdb.Ping(context.TODO()).Result()
	if err != nil {
		log.Fatalln(err)
	}
}

