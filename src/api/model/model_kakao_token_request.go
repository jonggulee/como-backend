package model

type KakaoTokenRequest struct {
	// authorization_code로 고정
	GrantType string `json:"grantType,omitempty"`
	// 앱 REST API 키
	ClientId string `json:"clientId,omitempty"`
	// 인가 코드가 리다이렉트된 URI
	RedirectURL string `json:"userId,omitempty"`
	// 인가 코드 받기 요청으로 얻은 인가 코드
	Code string `json:"code,omitempty"`
	// 앱 시크릿 키
	ClientSecret string `json:"clientSecret,omitempty"`
}
