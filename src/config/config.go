package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

type Cfg struct {
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		Db       string
	}
	ServerPort    int
	MosquittoHost string
}

var conf Cfg

func Init() {
	defer func() {
		if r := recover(); r != nil {
			logrus.Fatal("Failed to initialise configuration", zap.Any("error", r))
		}
	}()
	conf.Database.Host = readField("PG_HOST")
	conf.Database.Port = readField("PG_PORT")
	conf.Database.User = readField("PG_USER")
	conf.Database.Password = readField("PG_PASSWORD")
	conf.Database.Db = readField("PG_DATABASE")
	conf.ServerPort = readIntField("PORT")
	conf.MosquittoHost = readField("MQTT_HOST")
}

func Config() Cfg {
	return conf
}

func readField(f string) string {
	val := os.Getenv(f)
	if val == "" {
		panic(fmt.Errorf("key %s is not present. panic", f))
	}
	return val
}

func readIntField(f string) int {
	val := os.Getenv(f)
	if val == "" {
		panic(fmt.Errorf("key %s is not present. panic", f))
	}
	res, err := strconv.Atoi(val)
	if err != nil {
		panic(fmt.Errorf("key value can not be parsed, key: %s", f))
	}
	return res
}
