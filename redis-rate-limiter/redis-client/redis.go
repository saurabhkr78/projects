package redisclient

import (
	"github.com/redis/go-redis/v9"
)

type Client struct {
	//redis client instance
	RedisClient *redis.Client
}

func NewClient() *Client {
	return &Client{
		RedisClient: redis.NewClient(&redis.Options{
			Addr: "localhost:6379"}),
	}
}
