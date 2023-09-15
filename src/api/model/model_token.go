package model

import "time"

type KakaoToken struct {
	// Bearer 토큰
	TokenType string `json:"tokenType,omitempty"`
	// 액세스 토큰
	Token string `json:"token,omitempty"`
	// 액세스 토큰 만료 시간(초)
	Expiry time.Time `json:"expiresIn,omitempty"`
	// 리프레시 토큰
	RefreshToken string `json:"refreshToken,omitempty"`
	// 리프레시 토큰 만료 시간(초)
	RefreshTokenExpiresIn int `json:"refreshTokenExpiresIn,omitempty"`
	// 액세스 토큰으로 요청할 수 있는 권한 범위
	Scope string `json:"scope,omitempty"`
}
