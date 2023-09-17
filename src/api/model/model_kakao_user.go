package model

type KakaoUser struct {
	// 사용자 고유 ID
	Id int64 `json:"id,omitempty"`
	// 사용자의 닉네임
	Nickname string `json:"nickname,omitempty"`
	// 사용자가 앱에 등록한 이메일 주소
	Email string `json:"email,omitempty"`
}

type TempKakaoUser struct {
	ID         int64 `json:"id"`
	Properties struct {
		Nickname string `json:"nickname"`
	} `json:"properties"`
	KakaoAccount struct {
		Email string `json:"email"`
	} `json:"kakao_account"`
}
