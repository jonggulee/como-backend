package model

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type User struct {
	// 사용자 고유 ID
	Id int `json:"id,omitempty" gorm:"column:id;primaryKey;not null;autoIncrement:true"`

	KakaoId int64 `json:"kakaoId,omitempty" gorm:"column:kakao_id;unique;"`

	// 사용자가 앱에 등록한 이메일 주소
	Email string `json:"email,omitempty" gorm:"column:email;unique;not null"`

	// 사용자의 닉네임
	Nickname string `json:"nickname,omitempty" gorm:"column:nickname;not null"`

	// 가입 형태
	// 1: 카카오 계정으로 가입
	// 2: 이메일로 가입
	JoinedType int `json:"joinedType,omitempty" gorm:"column:joined_type;not null"`

	// 회원 가입 일시
	CreatedAt time.Time `json:"createdAt,omitempty" gorm:"column:created_at;not null;autoCreateTime"`

	// 회원 정보 수정 일시
	UpdatedAt time.Time `json:"updatedAt,omitempty" gorm:"column:updated_at;not null;autoUpdateTime:true"`

	// 회원 탈퇴 여부
	DeletedYn byte `json:"deleteYn,omitempty" gorm:"column:deleted_yn;not null"`
}

func (u *User) MarshalJSON() ([]byte, error) {
	type Alias User
	return json.Marshal(&struct {
		*Alias
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
	}{
		Alias:     (*Alias)(u),
		CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: u.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

func (u *User) AfterUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) TableName() string {
	return "user"
}
