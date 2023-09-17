package model

type KakaoUser struct {
	ID         int64 `json:"id"`
	Properties struct {
		Nickname string `json:"nickname"`
	} `json:"properties"`
	KakaoAccount struct {
		Email string `json:"email"`
	} `json:"kakao_account"`
}
