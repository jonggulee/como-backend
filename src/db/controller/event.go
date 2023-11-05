package controller

import (
	"github.com/jonggulee/como-backend/src/api/model"
	"github.com/jonggulee/como-backend/src/logger"
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

func EventSelectById(db *gorm.DB, reqId string, eventId int) error {
	logger.Debugf(reqId, "Try to select event by id %d", eventId)

	event := model.Event{}

	result := db.Table("event").Where("id = ? AND deleted_yn = 0", eventId).First(&event)
	if result.Error != nil {
		logger.Errorf(reqId, "Failed to select event by id %d, %s", eventId, result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		logger.Errorf(reqId, "No event found with id %d", eventId)
		return result.Error
	}

	return nil
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

	result := db.Table("event").Where("id = ? AND deleted_yn = 0", eventId).Updates(event)
	if result.Error != nil {
		logger.Errorf(reqId, "Failed to update event set ... %v", event)
		return result.Error
	}
	logger.Debugf(reqId, "rows affected: %d", result.RowsAffected)

	return nil
}

func EventDelete(db *gorm.DB, reqId string, eventId int, updateUserId int) error {
	logger.Debugf(reqId, "Try to delete from event where id = %d", eventId)

	updates := map[string]interface{}{
		"deleted_yn":     1,
		"update_user_id": updateUserId,
	}

	result := db.Table("event").Where("id = ?", eventId).Updates(updates)
	if result.Error != nil {
		logger.Errorf(reqId, "Failed to delete from event where id = %d", eventId)
		return result.Error
	}

	return nil
}
