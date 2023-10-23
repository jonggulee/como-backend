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

func EventCreateInsert(db *gorm.DB, reqId string, event *model.EventCreateRequest) error {
	logger.Debugf(reqId, "Try to insert into event ... values %v", event)

	result := db.Table("event").Create(event)
	if result.Error != nil {
		logger.Errorf(reqId, "Failed to insert into event... %s", result.Error)
		return result.Error
	}

	return nil
}

func EventEditPost(db *gorm.DB, reqId string, eventId int, event *model.EventCreateRequest) error {
	logger.Debugf(reqId, "Try to update event set ... %v", event)

	result := db.Table("event").Where("id = ?", eventId).Updates(event)
	if result.Error != nil {
		logger.Errorf(reqId, "Failed to update event set ... %v", event)
		return result.Error
	}

	return nil
}
