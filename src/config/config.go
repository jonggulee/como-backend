package config

type Config struct {
	ListenPort int

	// OAuth Kakao Client ID
	KakaoClientId string
	// OAuth Kakao Client Secret
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

type authSession struct {
	// 쿠키의 유효 시간
	MaxAge int
	// 쿠키의 도메인
	Path string
}

type kakaoEndpoint struct {
	AuthURL  string
	TokenURL string
	UserURL  string
}

var KakaoAuthSession = authSession{
	Path:   "/v1/user/login/kakao/token",
	MaxAge: 300,
}

var KakaoEndpoint = kakaoEndpoint{
	AuthURL:  "https://kauth.kakao.com/oauth/authorize",
	TokenURL: "https://kauth.kakao.com/oauth/token",
	UserURL:  "https://kapi.kakao.com/v2/user/me",
}
