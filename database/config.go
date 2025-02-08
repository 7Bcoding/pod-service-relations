package database

import (
	"github.com/spf13/viper"
)

const configKey = "datastore.mysql"

type config struct {
	URI      string
	MaxConn  int
	MaxIdle  int
	Lifetime int
}

func newDBConf() (config, error) {
	cfg := config{}
	err := viper.UnmarshalKey(configKey, &cfg)
	return cfg, err
}
