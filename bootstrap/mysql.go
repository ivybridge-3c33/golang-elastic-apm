package bootstrap

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ivybridge-3c33/golang-elastic-apm/configs"
	apmmysql "go.elastic.co/apm/module/apmgormv2/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type (
	// MySQL mysql database management
	MySQL struct {
	}
)

// dbMySQL variable for define connection
var dbMySQL map[string]*gorm.DB = make(map[string]*gorm.DB)

func CreateMySQLConnection(conf *configs.MySQLConn) *gorm.DB {
	connectionName := "default"
	if conf.ConnectionName != "" {
		connectionName = conf.ConnectionName
	}

	db, err := gorm.Open(apmmysql.Open(conf.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		panic(fmt.Sprintf("[mysql] failed to connect database: %s", err))
	}
	fmt.Println("[mysql] connected")

	if c, err := db.DB(); err != nil {
		panic(fmt.Sprintf("[mysql] connection poll error: %s", err))
	} else {
		if v := conf.MaxIdleConns; v > 0 {
			c.SetMaxIdleConns(v)
		}
		if v := conf.MaxOpenConns; v > 0 {
			c.SetMaxOpenConns(v)
		}
		if v := conf.ConnMaxIdleTime; v != nil {
			c.SetConnMaxIdleTime(*v)
		}
		if v := conf.ConnMaxLifetime; v != nil {
			c.SetConnMaxLifetime(*v)
		}
	}
	if debug, err := strconv.ParseBool(os.Getenv("APP_DEBUG")); err == nil {
		if debug {
			db = db.Debug()
		}
	}
	dbMySQL[connectionName] = db
	return db
}

// DB get mysql connection
func (c *MySQL) DB(connectionNames ...string) *gorm.DB {
	connectionName := "default"
	if len(connectionNames) > 0 {
		connectionName = connectionNames[0]
	}
	return dbMySQL[connectionName]
}

func (c *MySQL) HealthCheck(connectionNames ...string) string {
	connectionName := "default"
	if len(connectionNames) > 0 {
		connectionName = connectionNames[0]
	}
	err := dbMySQL[connectionName].Exec("SELECT 1").Error
	if err != nil {
		return "error " + err.Error()
	} else {
		return "ok"
	}
}
