package controller

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

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
	// sessionStore = sessions.NewCookieStore[any](sessions.DebugCookieConfig, []byte(sessionSecret), nil)
)

type KakaoTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

const (
	localServer = "http://localhost:8080"
)

func ReadKakaoConfig(cfg *config.Config) {
	oauthConf = &oauth2.Config{
		ClientID:     config.AppCtx.Cfg.KakaoClientId,
		ClientSecret: config.AppCtx.Cfg.KakaoClientSecret,

		RedirectURL: localServer + "/v1/user/login/kakao",
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
	state := authSession.Values["state"]
	fmt.Println("state: ", state)
	if state == nil {
		logger.Debugf(reqId, "state is %s", state)
		resp := newResponse(w, reqId, 400, "Bad Request")
		writeResponse(reqId, w, resp)
		return
	}
	if state != r.URL.Query().Get("state") {
		logger.Debugf(reqId, "state is %s", state)
		resp := newResponse(w, reqId, 403, "Forbidden")
		writeResponse(reqId, w, resp)
		return
	}

	// 쿠키에 저장된 state 보관 기간 만료
	authSession.Options = &sessions.Options{
		MaxAge: -1,
	}
	authSession.Save(r, w)

	// Get Kakao Auth Code
	c := r.URL.Query().Get("code")

	data := url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {oauthConf.ClientID},
		"redirect_uri":  {oauthConf.RedirectURL},
		"code":          {c},
		"client_secret": {oauthConf.ClientSecret},
	}

	// 아래 코드 수정 필요
	// fmt.Println(conf.Endpoint.TokenURL)
	resp, err := http.PostForm(oauthConf.Endpoint.TokenURL, data)
	if err != nil {
		fmt.Println("Error while posting:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error while reading response:", err)
		return
	}

	fmt.Println(string(body))

	// Unmarshal the JSON response into a KakaoTokenResponse struct
	var tokenResponse KakaoTokenResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		fmt.Println("Error while unmarshalling:", err)
		return
	}

	loginToken := &model.Token{}
	loginToken.Token = tokenResponse.AccessToken

	// resp = newOkResponse(w, reqId, constants.BASICOK)

	// resp :=

	// writeResponse(reqId, w, &model.Response{LoginToken: loginToken})

	// fmt.Println("test==============", loginToken.Token)

	// Output the token information
	// fmt.Printf("Access Token: %s\n", tokenResponse.AccessToken)
	// fmt.Printf("Expires In: %d\n", tokenResponse.ExpiresIn)
	// fmt.Printf("Refresh Token: %s\n", tokenResponse.RefreshToken)
	// fmt.Printf("Scope: %s\n", tokenResponse.Scope)
	// fmt.Printf("Token Type: %s\n", tokenResponse.TokenType)
}
