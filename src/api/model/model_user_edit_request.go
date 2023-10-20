package model

type UserEditRequest struct {
	// 사용자 ID
	UserId int `json:"userId,omitempty" gorm:"column:id;not null;index"`

	// 사용자 닉네임
	Nickname string `json:"nickname,omitempty" gorm:"column:nickname;not null"`
}
