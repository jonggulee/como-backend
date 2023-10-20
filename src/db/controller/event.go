package controller

import (
	"github.com/jonggulee/go-login.git/src/api/model"
	"github.com/jonggulee/go-login.git/src/logger"
	"gorm.io/gorm"
)

func EventSelect(db *gorm.DB, reqId string) (*[]model.Event, error) {
	logger.Debugf(reqId, "Try to select * from event")

	events := []model.Event{}

	result := db.Table("event").Find(&events, "deleted_yn=0")
	if result.Error != nil {
		logger.Errorf(reqId, "Failed to select * from event, %s", result.Error)
		return nil, result.Error
	}

	return &events, nil
}
