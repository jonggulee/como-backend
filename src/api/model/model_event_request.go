package model

type EventRequest struct {
	// 이벤트 ID
	Id int `json:"id,omitempty" gorm:"column:id;primaryKey;not null;autoIncrement:true"`

	// 이벤트 등록자 ID
	CreatedUserId int `json:"createdUserId,omitempty" gorm:"column:created_user_id;not null"`

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
}
