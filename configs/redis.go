package configs

import "github.com/go-redis/redis/v8"

type RedisConn struct {
	redis.UniversalOptions
	ConnectionName string //empty is default
}
