package model

import (
	"github.com/jonggulee/go-login.git/src/constants"
	"github.com/jonggulee/go-login.git/src/logger"
	"gorm.io/gorm"
)

type SessionRequest struct {
	UserId int `json:"userId,omitempty" gorm:"column:user_id;not null;index"`

	Session string `json:"session,omitempty" gorm:"column:session;not null'"`

	UserEmail string `json:"userEmail,omitempty" gorm:"column:user_email;not null;index"`
}

func (s *SessionRequest) BeforeCreate(tx *gorm.DB) error {
	logger.Debugf(constants.NOREQID, "Try to check duplication")

	foundSession := &SessionRequest{}
	result := tx.Table("session").Where("user_id = ? AND session = ? AND deleted_yn = 0", s.UserId, s.Session).Find(&foundSession)
	if result.Error != nil {
		logger.Errorf(constants.NOREQID, "Failed to select `session`, due to %s", result.Error)
		return result.Error
	}

	if s.UserId == foundSession.UserId && s.Session == foundSession.Session {
		logger.Errorf(constants.NOREQID, "Duplicated session")
		return result.Error
	}

	return nil
}
