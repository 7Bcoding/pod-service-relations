package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

const (
	configName = "config"
	//configName = "test_config"
	configType = "json"
	configPath = "."
)

// InitConfigs initializes the config from conf file and environment variables
func InitConfigs() {
	// set config name and location
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath(configPath)
	// set defaults
	setDefaults()
	// read config from specified file
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("read config file failed, using default values or env")
	}
	// override configs with environment variable
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	fmt.Println(viper.AllSettings())
	// auto-update config values
	// viper.WatchConfig()
}
func setDefaults() {
	viper.SetDefault("server.int.bind", "127.0.0.1")
	viper.SetDefault("server.int.port", 8088)
	viper.SetDefault("server.vpc.bind", "127.0.0.1")
	viper.SetDefault("server.vpc.port", 8088)
}
