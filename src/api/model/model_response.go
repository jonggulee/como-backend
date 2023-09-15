package model

type Response struct {
	RequestId     string `json:"requestId,omitempty"`
	StatusMessage string `json:"statusMessage,omitempty"`
	StatusCode    int    `json:"statusCode,omitempty"`

	KakaoToken   *Token `json:"loginToken,omitempty"`
	KakaoAuthUrl string `json:"kakaoAuthUrl,omitempty"`
}
