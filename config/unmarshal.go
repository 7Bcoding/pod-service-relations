package config

import (
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

// UnmarshalKey is a replacement of viper.UnmarshalKey, which cannot unmarshal the environment variables or the pflags.
func UnmarshalKey(conf *viper.Viper, key string, container interface{}) {
	// panics if container is not a pointer
	t := reflect.TypeOf(container).Elem()
	v := reflect.ValueOf(container).Elem()
	dfs(t, v, key, conf)
}

func setValue(name string, v reflect.Value, conf *viper.Viper) {
	if strings.HasPrefix(name, ".") {
		name = name[1:]
	}
	switch v.Kind() {
	case reflect.Bool:
		v.SetBool(conf.GetBool(name))
	case reflect.String:
		v.SetString(conf.GetString(name))
	case reflect.Int:
		v.SetInt(int64(conf.GetInt(name)))
	case reflect.Int64:
		v.SetInt(conf.GetInt64(name))
	case reflect.Float64:
		v.SetFloat(float64(conf.GetFloat64(name)))
	default:
		// non-recoverable situation
		panic("invalid type")
	}
}

func dfs(t reflect.Type, v reflect.Value, name string, conf *viper.Viper) {
	k := v.Kind()
	if k != reflect.Struct {
		setValue(name, v, conf)
		return
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		dfs(field.Type, v.Field(i), name+"."+field.Name, conf)
	}
}
