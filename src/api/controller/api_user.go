package controller

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/jonggulee/go-login.git/src/api/model"
	"github.com/jonggulee/go-login.git/src/config"
	"github.com/jonggulee/go-login.git/src/constants"
	"github.com/jonggulee/go-login.git/src/logger"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/kakao"
)

var (
	oauthConf *oauth2.Config
	store     = sessions.NewCookieStore([]byte("secret"))
)

const (
	localServer = "http://localhost:8080"
)

func ReadKakaoConfig(cfg *config.Config) {
	oauthConf = &oauth2.Config{
		ClientID:     config.AppCtx.Cfg.KakaoClientId,
		ClientSecret: config.AppCtx.Cfg.KakaoClientSecret,

		RedirectURL: localServer + config.KakaoAuthSession.Path,
		Endpoint:    kakao.Endpoint,
	}
}

func randomState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func KakaoAuthUrlGet(w http.ResponseWriter, r *http.Request) {
	reqId := getRequestId(w, r)
	logger.Debugf(reqId, "user/login/kakao/authurl GET started")

	authSession, _ := store.Get(r, "authSession")
	authSession.Options = &sessions.Options{
		Path:   config.KakaoAuthSession.Path,
		MaxAge: config.KakaoAuthSession.MaxAge,
	}

	state := randomState()
	authSession.Values["state"] = state
	authSession.Save(r, w)

	url := oauthConf.AuthCodeURL(state, oauth2.AccessTypeOffline)

	resp := newOkResponse(w, reqId, constants.BASICOK)
	resp.KakaoAuthUrl = url
	writeResponse(reqId, w, resp)
}

func KakaoTokenGet(w http.ResponseWriter, r *http.Request) {
	reqId := getRequestId(w, r)
	logger.Debugf(reqId, "user/login/kakao GET started")

	// 쿠키에 저장된 state 값 불러오기
	authSession, _ := store.Get(r, "authSession")
	s := authSession.Values["state"]
	if s == nil {
		logger.Debugf(reqId, "state is %s", s)
		resp := newResponse(w, reqId, 400, "Bad Request")
		writeResponse(reqId, w, resp)
		return
	}
	if s != r.URL.Query().Get("state") {
		logger.Debugf(reqId, "state is %s", s)
		resp := newResponse(w, reqId, 403, "Forbidden")
		writeResponse(reqId, w, resp)
		return
	}

	// 쿠키에 저장된 state 보관 기간 만료
	authSession.Options = &sessions.Options{
		MaxAge: -1,
	}
	authSession.Save(r, w)

	// 인가 코드 받기
	c := r.URL.Query().Get("code")

	ctx := r.Context()
	token, err := oauthConf.Exchange(ctx, c)
	if err != nil {
		logger.Errorf(reqId, "Failed to get token")
		resp := newResponse(w, reqId, 500, "Internal Server Error")
		writeResponse(reqId, w, resp)
		return
	}

	kakaoToken := &model.KakaoToken{}
	kakaoToken.TokenType = token.TokenType
	kakaoToken.Token = token.AccessToken
	kakaoToken.RefreshToken = token.RefreshToken
	kakaoToken.Expiry = token.Expiry

	resp := newOkResponse(w, reqId, constants.BASICOK)
	resp.KakaoToken = kakaoToken
	writeResponse(reqId, w, resp)
}
