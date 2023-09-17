package controller

import (
	"github.com/jonggulee/go-login.git/src/api/model"
	"github.com/jonggulee/go-login.git/src/logger"
	"gorm.io/gorm"
)

func UserSignUp(db *gorm.DB, reqId string, user *model.UserInfo) error {
	logger.Debugf(reqId, "Try to insert into user ... values %v", user)

	result := db.Table("user").Create(user)
	if result.Error != nil {
		logger.Errorf(reqId, "Failed to insert into user... %s", result.Error)
		return result.Error
	}

	return nil
}

func UserFindByKakaoIdSelect(db *gorm.DB, reqId string, user *model.UserInfo) (*model.UserInfo, error) {
	logger.Debugf(reqId, "Try to select from user where kakao_id = %v", user.KakaoId)

	result := db.Table("user").
		Where("kakao_id = ?", user.KakaoId).
		Find(&user, "deleted_yn=0")

	if result.RowsAffected == 0 {
		return nil, nil
	}

	if result.Error != nil {
		logger.Errorf(reqId, "Failed to select from user... %s", result.Error)
		return nil, result.Error
	}

	return user, nil
}
