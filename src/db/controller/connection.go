package controller

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jonggulee/como-backend/src/constants"
	"github.com/jonggulee/como-backend/src/db/model"
	"github.com/jonggulee/como-backend/src/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func Connection(dbName, user, password, address string, port int) (*model.DbConnection, error) {
	logger.Debugf(constants.NOREQID, "Try to connect %s:%d/%s", address, port, dbName)

	options := "?charset=utf8mb4&parseTime=True"
	newLogger := gormLogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		gormLogger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  gormLogger.Error,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
	configs := &gorm.Config{
		Logger: newLogger,
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s%s", user, password, address, port, dbName, options)
	db, err := gorm.Open(mysql.Open(dsn), configs)
	if err != nil {
		logger.Errorf(constants.NOREQID, "Failed to connect DB(%s:%d): %s", address, port, err)
		return nil, err
	}

	logger.Debugf(constants.NOREQID, "Connected on %s:%d/%s", address, port, dbName)

	return &model.DbConnection{
		Db: db,
	}, nil
}
