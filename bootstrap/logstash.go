package bootstrap

import (
	"fmt"
	"net"
	"os"

	logrustash "github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/sirupsen/logrus"
	"go.elastic.co/apm/module/apmlogrus"
)

var Logger *logrus.Logger

func CreateLogstashConnection(serviceName string) {
	Logger = logrus.New()
	Logger.AddHook(&apmlogrus.Hook{})
	conn, err := net.Dial("tcp", os.Getenv("LOGSTASH_HOST"))
	if err != nil {
		panic("[Logstash] error : " + err.Error())
	}
	fmt.Println("[Logstash] connected")
	hook := logrustash.New(conn, logrustash.DefaultFormatter(logrus.Fields{
		"service_name": serviceName,
		"app_env":      os.Getenv("APP_ENV"),
		"app_debug":    os.Getenv("APP_DEBUG"),
		"app_version":  os.Getenv("APP_VERSION"),
	}))

	Logger.SetReportCaller(true)
	Logger.Hooks.Add(hook)
}
