package client

import (
	"github.com/redis/go-redis/v9"
)

type Client struct {
	//redis client instance
	redisClient *redis.Client
}

func NewClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379"})
}
