package controller

import (
	"github.com/jonggulee/como-backend/src/api/model"
	"github.com/jonggulee/como-backend/src/logger"
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

func SessionSelectByUserId(db *gorm.DB, reqId string, userId int) (*model.Session, error) {
	logger.Debugf(reqId, "Try to select * from session where user_id = %d", userId)

	sessionInfo := &model.Session{}

	result := db.Table("session").
		Where("session.user_id = ? AND session.deleted_yn = 0", userId).
		Find(sessionInfo)

	if result.Error != nil {
		logger.Errorf(reqId, "Failed to select * from session where user_id = %d", userId)
		return nil, result.Error
	}

	return sessionInfo, nil
}

func SessionDelete(db *gorm.DB, reqId string, sessionId string) error {
	logger.Debugf(reqId, "Try to delete from session where session = %s", sessionId)

	result := db.Table("session").Where("session = ?", sessionId).Update("deleted_yn", 1)
	if result.Error != nil {
		logger.Errorf(reqId, "Failed to delete from session where session = %s", sessionId)
		return result.Error
	}

	return nil
}

func UserDetailSelect(db *gorm.DB, reqId string, userId int) (*model.User, error) {
	logger.Debugf(reqId, "Try to select * from user where id = %d", userId)

	user := model.User{}

	result := db.Table("user").Find(&user, "id = ? AND deleted_yn = 0", userId)

	if result.RowsAffected == 0 {
		return nil, nil
	}

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

func UserReactivateByEmail(db *gorm.DB, reqId string, user *model.User) error {
	logger.Debugf(reqId, "Try to update user set ... %v", user)

	result := db.Table("user").Omit("UpdatedAt").Where("email = ? AND deleted_yn = 1", user.Email).UpdateColumn("deleted_yn", 0)
	if result.Error != nil {
		logger.Errorf(reqId, "Failed to delete user set ... %v, %s", user, result.Error)
		return result.Error
	}

	return nil
}

func UserWithdrawDelete(db *gorm.DB, reqId string, user *model.User) error {
	logger.Debugf(reqId, "Try to delete user set ... %v", user)

	result := db.Table("user").Omit("UpdatedAt").Where("deleted_yn = 0").Updates(user)
	if result.Error != nil {
		logger.Errorf(reqId, "Failed to update delete user set ... %v, %s", user, result.Error)
		return result.Error
	}

	return nil
}

func UserFindByKakaoIdSelect(db *gorm.DB, reqId string, user *model.User) (*model.User, error) {
	logger.Debugf(reqId, "Try to select from user where kakao_id = %v", user.KakaoId)

	result := db.Table("user").Where("kakao_id = ?", user.KakaoId).Find(&user, "deleted_yn = 0")

	if result.RowsAffected == 0 {
		return nil, nil
	}

	if result.Error != nil {
		logger.Errorf(reqId, "Failed to select from user... %s", result.Error)
		return nil, result.Error
	}

	return user, nil
}
