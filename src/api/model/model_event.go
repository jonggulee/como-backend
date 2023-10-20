package model

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type Event struct {
	// 이벤트 고유 ID
	Id int `json:"id,omitempty" gorm:"column:id;primaryKey;not null;autoIncrement:true"`

	// 이벤트 제목
	Title string `json:"title,omitempty" gorm:"column:title;not null"`

	// 이벤트 내용
	Content string `json:"content,omitempty" gorm:"column:content;not null"`

	// 이벤트 이미지 URL
	ImageUrl string `json:"imageUrl,omitempty" gorm:"column:image_url;not null"`

	// 이벤트 시작 일시
	StartDate string `json:"startDate,omitempty" gorm:"column:start_date;not null"`

	// 이벤트 종료 일시
	EndDate string `json:"endDate,omitempty" gorm:"column:end_date;not null"`

	// 이벤트 등록자 ID
	CreatedUserId int `json:"createdUserId,omitempty" gorm:"column:created_user_id;not null"`

	// 이벤트 수정자 ID
	UpdatedUserId int `json:"updatedUserId,omitempty" gorm:"column:updated_user_id;not null"`

	// 이벤트 등록 일시
	CreatedAt time.Time `json:"createdAt,omitempty" gorm:"column:created_at;not null;autoCreateTime"`

	// 이벤트 수정 일시
	UpdatedAt time.Time `json:"updatedAt,omitempty" gorm:"column:updated_at;not null;autoUpdateTime:true"`

	// 이벤트 삭제 여부
	DeletedYn byte `json:"deleteYn,omitempty" gorm:"column:deleted_yn;not null"`
}

func (e *Event) MarshalJSON() ([]byte, error) {
	type Alias Event
	return json.Marshal(&struct {
		*Alias
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
	}{
		Alias:     (*Alias)(e),
		CreatedAt: e.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: e.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

func (e *Event) AfterUpdate(tx *gorm.DB) error {
	e.UpdatedAt = time.Now()
	return nil
}

func (e *Event) TableName() string {
	return "event"
}
