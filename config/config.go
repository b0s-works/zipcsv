package config

import (
	"github.com/sirupsen/logrus"
	"github.com/jinzhu/configor"
)

type Config struct {
	Aggregation map[string]string
}

func New() Config {
	var config Config
	configor.Load(&config, "config/config.json")
	logrus.Debugf("config: %#v", config)

	return config
}
