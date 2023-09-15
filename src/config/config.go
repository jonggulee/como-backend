package config

type Config struct {
	ListenPort int

	KakaoClientId     string
	KakaoClientSecret string
}

type AuthSession struct {
	// 쿠키의 유효 시간
	MaxAge int
	// 쿠키의 도메인
	Path string
}

var KakaoAuthSession = AuthSession{
	Path:   "/v1/user/login/kakao",
	MaxAge: 300,
}
