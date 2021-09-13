package bootstrap

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/ivybridge-3c33/golang-elastic-apm/configs"
	apmgoredis "go.elastic.co/apm/module/apmgoredisv8"
)

type (
	// RedisDB database management
	RedisDB struct {
	}
)

// dbRedis variable for define connection
var dbRedis map[string]redis.UniversalClient = make(map[string]redis.UniversalClient)

// CreateRedisConnection make connection
// example
// database, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
// bootstrap.CreateRedisConnection(&redis.UniversalOptions{
//  Addrs:       strings.Split(os.Getenv("REDIS_HOST"), ","),
//  Password:    os.Getenv("REDIS_PASSWORD"),
//  DB:          database,
//  DialTimeout: time.Duration(15) * time.Second,
//  ConnectionName: "test" // empty is default
// })
// new(bootstrap).DB().Get()....
// new(bootstrap).DB("test").Get()....
func CreateRedisConnection(conf *configs.RedisConn) redis.UniversalClient {
	connectionName := "default"
	if conf.ConnectionName != "" {
		connectionName = conf.ConnectionName
	}
	db := redis.NewUniversalClient(&conf.UniversalOptions)
	db.AddHook(apmgoredis.NewHook())
	if _, err := db.Ping(context.TODO()).Result(); err != nil {
		panic(fmt.Sprintf("[redis] connect database fail as error: %s", err))
	}
	fmt.Printf("[redis] connected\n")
	dbRedis[connectionName] = db
	return db
}

// DB get redis connection
func (c *RedisDB) DB(connectionNames ...string) redis.UniversalClient {
	connectionName := "default"
	if len(connectionNames) > 0 {
		connectionName = connectionNames[0]
	}
	return dbRedis[connectionName]
}

func (c *RedisDB) HealthCheck(connectionNames ...string) string {
	connectionName := "default"
	if len(connectionNames) > 0 {
		connectionName = connectionNames[0]
	}
	_, err := dbRedis[connectionName].Ping(context.TODO()).Result()
	if err != nil {
		return "error " + err.Error()
	} else {
		return "ok"
	}
}
