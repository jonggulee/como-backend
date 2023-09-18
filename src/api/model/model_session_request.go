package model

type SessionRequest struct {
	UserId int `json:"userId,omitempty" gorm:"column:user_id;not null;index"`

	Session string `json:"session,omitempty" gorm:"column:session;not null'"`

	UserEmail string `json:"userEmail,omitempty" gorm:"column:user_email;not null;index"`
}
