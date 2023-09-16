package config

type Config struct {
	ListenPort int

	// OAuth kakao Client ID
	KakaoClientId string
	// OAuth kakao Client Secret
	KakaoClientSecret string

	// DB Address
	DbAddress string
	// DB Port
	DbPort int
	// DB Name
	DbName string
	// DB User
	DbUser string
	// DB Password
	DbPassword string
}

type AuthSession struct {
	// 쿠키의 유효 시간
	MaxAge int
	// 쿠키의 도메인
	Path string
}

var KakaoAuthSession = AuthSession{
	Path:   "/v1/user/login/kakao/token",
	MaxAge: 300,
}
