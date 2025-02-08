package config

import (
	"encoding/json"
	"github.com/spf13/viper"
)

const (
	serverConfigKey        = "server"
	filecenterConfigKey    = "filecenter"
	irepoConfigKey         = "irepo"
	hanoiConfigKey         = "hanoi"
	tempDirConfigKey       = "tempdir"
	KubeApiServerConfigKey = "kube_api_server"
)

type ServerConfig struct {
	Bind    string
	Port    int
	Session struct {
		Name   string
		Secret string
		//Redis  struct {
		//	URI     string
		//	MaxIdle int
		//}
	}
}

type FilecenterConfig struct {
	UploadAPI string
	FileName  string
	Token     string
	Group     string
	Password  string
	Timeout   int
}

type IrepoConfig struct {
	Plat      string
	PlatToken string
	Timeout   int
}

type HanoiConfig struct {
	User    string
	Token   string
	Timeout int
}

type TempDIRConfig struct {
	Dir    string
	Prefix string
}

type KubeApiServerConfig struct {
	KubeConfigPath string `json:"kube_config_path"`
}

func NewServerConfig() ServerConfig {
	conf := ServerConfig{}
	//UnmarshalKey(viper.GetViper(), serverConfigKey, &conf)
	//return conf
	serverConfigData := viper.Get(serverConfigKey)
	arr, err := json.Marshal(serverConfigData)
	if err != nil {
		panic(err)
	}
	// 反序列化
	err2 := json.Unmarshal(arr, &conf)
	if err2 != nil {
		panic(err2)
	}
	return conf
}

func NewFilecenterConfig() FilecenterConfig {
	conf := FilecenterConfig{}
	UnmarshalKey(viper.GetViper(), filecenterConfigKey, &conf)
	return conf
}

func NewIrepoConfig() IrepoConfig {
	conf := IrepoConfig{}
	UnmarshalKey(viper.GetViper(), irepoConfigKey, &conf)
	return conf
}

func NewHanoiConfig() HanoiConfig {
	conf := HanoiConfig{}
	UnmarshalKey(viper.GetViper(), hanoiConfigKey, &conf)
	return conf
}

func NewTempDirConfig() TempDIRConfig {
	conf := TempDIRConfig{}
	UnmarshalKey(viper.GetViper(), tempDirConfigKey, &conf)
	return conf
}

func NewKubeApiServerConfig() KubeApiServerConfig {
	conf := KubeApiServerConfig{}
	KubeApiServerConfigData := viper.Get(KubeApiServerConfigKey)
	arr, err := json.Marshal(KubeApiServerConfigData)
	if err != nil {
		panic(err)
	}
	// 反序列化
	err2 := json.Unmarshal(arr, &conf)
	if err2 != nil {
		panic(err2)
	}
	return conf
}
