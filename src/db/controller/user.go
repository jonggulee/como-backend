package controller

import (
	"github.com/jonggulee/go-login.git/src/api/model"
	"github.com/jonggulee/go-login.git/src/logger"
	"gorm.io/gorm"
)

func SessionInsert(db *gorm.DB, reqId string, session *model.SessionRequest) error {
	logger.Debugf(reqId, "Try to insert into session ... values %v", session)

	result := db.Table("session").Create(session)
	if result.Error != nil {
		logger.Errorf(reqId, "Failed to insert into session... %s", result.Error)
		return result.Error
	}

	return nil
}

func SessionSelect(db *gorm.DB, reqId, session string) (*model.Session, error) {
	logger.Debugf(reqId, "Try to select * from session where session = %s", session)

	sessionInfo := &model.Session{}

	result := db.Table("session").
		Joins("JOIN user ON session.user_id = user.id").
		Where("session.session = ? AND session.deleted_yn = 0", session).
		Find(sessionInfo)

	if result.Error != nil {
		logger.Errorf(reqId, "Failed to select * from session where session = %s", session)
		return nil, result.Error
	}

	return sessionInfo, nil
}

func UserDetailSelect(db *gorm.DB, reqId string, userId int) (*model.User, error) {
	logger.Debugf(reqId, "Try to select * from user where id = %d", userId)

	user := model.User{}

	result := db.Table("user").Find(&user, "id = ? AND deleted_yn = 0", userId)
	if result.Error != nil {
		logger.Errorf(reqId, "Failed to select * from user where id = %d", userId)
		return nil, result.Error
	}

	return &user, nil
}

func UserDetailUpdate(db *gorm.DB, reqId string, userReq *model.UserEditRequest) error {
	logger.Debugf(reqId, "Try to update user set ... %v", userReq)

	result := db.Table("user").Omit("UpdatedAt").Where("deleted_yn = 0").Updates(userReq)
	if result.Error != nil {
		logger.Errorf(reqId, "Failed to update user set ... %v, %s", userReq, result.Error)
		return result.Error
	}

	return nil
}

func UserSignUp(db *gorm.DB, reqId string, user *model.User) error {
	logger.Debugf(reqId, "Try to insert into user ... values %v", user)

	result := db.Table("user").Create(user)
	if result.Error != nil {
		logger.Errorf(reqId, "Failed to insert into user... %s", result.Error)
		return result.Error
	}

	return nil
}

func UserFindByKakaoIdSelect(db *gorm.DB, reqId string, user *model.User) (*model.User, error) {
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
