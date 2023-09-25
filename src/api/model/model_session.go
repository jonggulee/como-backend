package model

import (
	"encoding/json"
	"time"
)

type Session struct {
	// 세션 ID
	Id int `json:"id,omitempty" gorm:"column:id;not null;primaryKey;index"`

	// 사용자 ID
	UserId int `json:"userId,omitempty" gorm:"column:user_id;not null;index"`

	// 사용자
	User *User `json:"user,omitempty" gorm:"foreignKey:UserId"`

	// 세션
	Session string `json:"session,omitempty" gorm:"column:session;not null'"`

	// 세션 생성 일시
	CreatedAt time.Time `json:"createdAt,omitempty" gorm:"column:created_at;not null;autoCreateTime"`

	// 세션 수정 일시
	UpdatedAt time.Time `json:"updatedAt,omitempty" gorm:"column:updated_at;not null;autoUpdateTime"`

	// 세션 삭제 시간
	DeletedYn byte `json:"-" gorm:"column:deleted_yn;not null;default:0"`
}

func (s *Session) MarshalJSON() ([]byte, error) {
	type Alias Session
	return json.Marshal(&struct {
		*Alias
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
	}{
		Alias:     (*Alias)(s),
		CreatedAt: s.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: s.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

func (s *Session) TableName() string {
	return "session"
}
