package redis

import (
	"github.com/redis/go-redis/v9"
)

// RedisEngine is an interface for Redis client.
type RedisEngine interface {
	GetRedisClient() *redis.Client
	Close()
}
