package database

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

// Init initializes the database
func Init() {
	var err error
	// load config
	conf, err := newDBConf()
	if err != nil {
		panic(err)
	}
	// log config
	ormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Warn, // Log level
			Colorful:      true,        // Disable color
		},
	)
	// connect
	db, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       conf.URI,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger: ormLogger,
	})
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	// Set connection pool
	sqlDB.SetMaxOpenConns(conf.MaxConn)
	sqlDB.SetMaxIdleConns(conf.MaxIdle)
	sqlDB.SetConnMaxLifetime(time.Duration(conf.Lifetime) * time.Second)
	// Disable pluralization, which is annoying
	//sqlDB.SingularTable(true)
	// `true` for detailed logs, `false` for no log, default, will only print error logs
	//sqlDB.LogMode(true)
}

// GetDB returns a gorm.DB
func GetDB() *gorm.DB {
	return db
}

// Close closes the database connections
func Close() {
	//db.Close()
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.Close()
}
