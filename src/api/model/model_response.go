package model

type Response struct {
	RequestId     string `json:"requestId,omitempty"`
	StatusMessage string `json:"statusMessage,omitempty"`
	StatusCode    int    `json:"statusCode,omitempty"`

	KakaoToken   *KakaoToken `json:"kakaoToken,omitempty"`
	KakaoUser    *KakaoUser  `json:"kakaoUser,omitempty"`
	KakaoAuthUrl string      `json:"kakaoAuthUrl,omitempty"`
}
