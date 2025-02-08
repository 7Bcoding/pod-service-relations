package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var levelMap = map[string]logrus.Level{
	"info":  logrus.InfoLevel,
	"error": logrus.ErrorLevel,
	"debug": logrus.DebugLevel,
	"warn":  logrus.WarnLevel,
	"trace": logrus.TraceLevel,
	"panic": logrus.PanicLevel,
	"fatal": logrus.FatalLevel,
}
var logger *logrus.Logger

// Init initializes the global logger
func Init() {
	logrus.SetReportCaller(true)
	logger = logrus.New()
	if _, err := os.Stat(viper.GetString("log.dir")); os.IsNotExist(err) {
		err := os.MkdirAll(viper.GetString("log.dir"), 0700)
		if err != nil {
			fmt.Println(err)
		}
	}
	file, err := os.OpenFile(
		filepath.Join(viper.GetString("log.dir"), viper.GetString("log.file")),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666,
	)
	if err != nil {
		panic(err)
	}
	logger.SetLevel(levelMap[viper.GetString("log.level")])
	logger.SetOutput(file)
	logger.SetReportCaller(true)
	logger.SetFormatter(&nested.Formatter{
		FieldsOrder:     []string{"component", "category"},
		TimestampFormat: time.RFC3339,
	})
}

// GetLogger returns the global logger
func GetLogger() *logrus.Logger {
	return logger
}
