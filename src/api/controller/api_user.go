package controller

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/jonggulee/go-login.git/src/api/model"
	"github.com/jonggulee/go-login.git/src/config"
	"github.com/jonggulee/go-login.git/src/constants"
	dbController "github.com/jonggulee/go-login.git/src/db/controller"
	"github.com/jonggulee/go-login.git/src/logger"
	"github.com/jonggulee/go-login.git/src/utils"
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

func accessTokenGet(w http.ResponseWriter, r *http.Request, user *model.User) (*model.Token, error) {
	reqId := getRequestId(w, r)
	logger.Debugf(reqId, "user/login POST started")

	// 세션 발급
	session := &model.SessionRequest{}
	userSession := uuid.New().String()
	session.UserId = user.Id
	session.UserEmail = user.Email
	session.Session = userSession

	err := dbController.SessionInsert(config.AppCtx.Db.Db, reqId, session)
	if err != nil {
		logger.Errorf(reqId, "Failed to insert into session ... values %v, duo to %s", session, err)
		return nil, err
	}

	// access token 발급
	accessToken := jwt.New(jwt.SigningMethodHS256)
	claims := accessToken.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(time.Hour * 168).Unix()
	claims["exp"] = time.Now().Add(time.Minute * 1).Unix()
	claims["iat"] = time.Now().Unix()
	claims["session"] = userSession

	t, err := accessToken.SignedString([]byte(constants.JWTPSK))
	if err != nil {
		logger.Errorf(reqId, "Failed to set jwt token, %s", err)
		return nil, err
	}
	if t == "" {
		logger.Errorf(reqId, "Failed to get jwt token")
		return nil, err
	}

	loginToken := &model.Token{}
	loginToken.Token = t

	return loginToken, nil
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

	accessToken, err := accessTokenGet(w, r, user)
	if err != nil {
		logger.Errorf(reqId, "Failed to login user %s", err)
		resp := newResponse(w, reqId, 500, "Internal Server Error")
		writeResponse(reqId, w, resp)
		return
	}

	resp := newOkResponse(w, reqId, constants.BASICOK)
	resp.LoginToken = accessToken
	writeResponse(reqId, w, resp)
}

func DetailUserGet(w http.ResponseWriter, r *http.Request) {
	reqId := getRequestId(w, r)
	logger.Debugf(reqId, "user/detail GET started")

	token := r.Header.Get("Authorization")
	if token == "" {
		logger.Errorf(reqId, "Failed to get authorization header")
		resp := newResponse(w, reqId, 400, "Bad Request")
		writeResponse(reqId, w, resp)
		return
	}

	decodedJwt, err := utils.DecodeJwt(reqId, token)
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&(jwt.ValidationErrorExpired) != 0 {
				logger.Errorf(reqId, "Failed to expired or not valid yet")
				resp := newResponse(w, reqId, 401, "Token Expired")
				writeResponse(reqId, w, resp)
				return
			}
		}

		logger.Errorf(reqId, "Failed to decode authorization header to jwt token %s", err)
		resp := newResponse(w, reqId, 500, "Internal error")
		writeResponse(reqId, w, resp)
		return
	}

	if decodedJwt.Session == "" {
		logger.Errorf(reqId, "Failed to get session from jwt token")
		resp := newResponse(w, reqId, 403, "Forbidden")
		writeResponse(reqId, w, resp)
		return
	}

	user, err := dbController.UserDetailSelect(config.AppCtx.Db.Db, reqId, decodedJwt.UserId)
	if err != nil {
		logger.Errorf(reqId, "Failed to select * from user where user_id = %d", decodedJwt.UserId)
		resp := newResponse(w, reqId, 500, "Internal Server Error")
		writeResponse(reqId, w, resp)
		return
	}
	if user == nil {
		logger.Errorf(reqId, "Failed to get user")
		resp := newResponse(w, reqId, 403, "Forbidden")
		writeResponse(reqId, w, resp)
		return
	}

	resp := newOkResponse(w, reqId, constants.BASICOK)
	resp.UserInfo = user
	writeResponse(reqId, w, resp)
}
