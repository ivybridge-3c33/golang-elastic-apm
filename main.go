package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/ivybridge-3c33/golang-elastic-apm/bootstrap"
	"github.com/ivybridge-3c33/golang-elastic-apm/configs"
	"github.com/ivybridge-3c33/golang-elastic-apm/entities"
	"github.com/ivybridge-3c33/golang-elastic-apm/handlers"
	"github.com/joho/godotenv"
	"go.elastic.co/apm/module/apmgin"
	"gorm.io/driver/mysql"
)

func init() {
	godotenv.Load()
}

func main() {
	serviceName := "test-service"
	bootstrap.CreateLogstashConnection(serviceName)
	maxConn := 2
	if v := os.Getenv("MYSQL_MAX_CONN"); v != "" {
		if newConn, err := strconv.Atoi(v); err == nil {
			maxConn = newConn
		}
	}
	mysqlDefaultConn := bootstrap.CreateMySQLConnection(&configs.MySQLConn{
		Config: mysql.Config{
			DSN: fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
				os.Getenv("MYSQL_USERNAME"),
				os.Getenv("MYSQL_PASSWORD"),
				os.Getenv("MYSQL_HOST"),
				os.Getenv("MYSQL_PORT"),
				os.Getenv("MYSQL_DBNAME"),
			),
		},
		MaxOpenConns: maxConn,
	})
	mysqlDB, _ := mysqlDefaultConn.DB()
	mysqlDefaultConn.AutoMigrate(&entities.Test{})
	defer mysqlDB.Close()

	rdDatabase, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	rdDB := bootstrap.CreateRedisConnection(&configs.RedisConn{
		UniversalOptions: redis.UniversalOptions{
			Addrs:    strings.Split(os.Getenv("REDIS_HOST"), ","),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       rdDatabase,
		},
	})
	defer rdDB.Close()

	engine := gin.New()
	engine.Use(apmgin.Middleware(engine))

	handler := new(handlers.TestHandler)
	engine.GET("/", handler.Home)
	engine.Run(":8088")
}
