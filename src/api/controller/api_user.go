package controller

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/jonggulee/go-login.git/src/api/model"
	"github.com/jonggulee/go-login.git/src/config"
	"github.com/jonggulee/go-login.git/src/constants"
	dbController "github.com/jonggulee/go-login.git/src/db/controller"
	"github.com/jonggulee/go-login.git/src/logger"
	"golang.org/x/oauth2"
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
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.KakaoEndpoint.AuthURL,
			TokenURL: config.KakaoEndpoint.TokenURL,
		},
	}
}

func randomState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func kakaoUserGet(w http.ResponseWriter, r *http.Request, token *model.KakaoToken) (*model.User, error) {
	reqId := getRequestId(w, r)
	logger.Debugf(reqId, "user/login/kakao/user GET started")

	url := config.KakaoEndpoint.UserURL

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Errorf(reqId, "Failed to create request %s", err)
		return nil, fmt.Errorf("failed to create request %s", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", token.TokenType, token.Token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf(reqId, "Failed to get user info %s", err)
		return nil, fmt.Errorf("failed to get user info %s", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("", "Failed to read body %s", err)
		return nil, fmt.Errorf("failed to read body %s", err)
	}

	var kakaoUser model.KakaoUser
	err = json.Unmarshal(body, &kakaoUser)
	if err != nil {
		logger.Errorf(reqId, "Failed to unmarshal body %s", err)
		return nil, fmt.Errorf("failed to unmarshal body %s", err)
	}

	user := &model.User{
		KakaoId:    kakaoUser.ID,
		Nickname:   kakaoUser.Properties.Nickname,
		Email:      kakaoUser.KakaoAccount.Email,
		JoinedType: 1,
	}

	return user, nil
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

	// kakao user 정보 가져오기
	kakaoUser, err := kakaoUserGet(w, r, kakaoToken)
	if err != nil {
		logger.Errorf(reqId, "Failed to get user info")
		resp := newResponse(w, reqId, 500, "Internal Server Error")
		writeResponse(reqId, w, resp)
		return
	}

	user, err := dbController.UserFindByKakaoIdSelect(config.AppCtx.Db.Db, reqId, kakaoUser)
	if err != nil {
		logger.Errorf(reqId, "Failed to select from user... %s", err)
		resp := newResponse(w, reqId, 500, "Internal Server Error")
		writeResponse(reqId, w, resp)
		return
	}

	if user == nil {
		// DB에 저장
		err = dbController.UserSignUp(config.AppCtx.Db.Db, reqId, kakaoUser)
		if err != nil {
			logger.Errorf(reqId, "Failed to insert into user ... values %v, duo to %s", kakaoUser, err)
			resp := newResponse(w, reqId, 500, "Internal Server Error")
			writeResponse(reqId, w, resp)
			return
		}
	}

	resp := newOkResponse(w, reqId, constants.BASICOK)
	resp.KakaoToken = kakaoToken
	writeResponse(reqId, w, resp)
}
