package model

type Session struct {
	// 세션 ID
	Id int `json:"id,omitempty" gorm:"column:id;not null;primaryKey;index"`

	// 사용자 ID
	UserId int `json:"-" gorm:"column:user_id;not null;index"`

	// 사용자
	User *User `json:"user,omitempty" gorm:"foreignKey:UserId"`
}
