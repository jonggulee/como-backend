package model

type LoginRequest struct {
	// 사용자 Email
	Email string `json:"email,omitempty"`
}
