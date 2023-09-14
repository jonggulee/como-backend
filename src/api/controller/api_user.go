package controller

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/sessions"
	"github.com/jonggulee/go-login.git/src/logger"
	"golang.org/x/oauth2"
)

const (
	localServer = "http://localhost:8080"
	authKakao   = "https://kauth.kakao.com"
)

var (
	store = sessions.NewCookieStore([]byte("secret"))

	conf = &oauth2.Config{
		// ClientID:     os.Getenv("CLIENT_ID"),
		ClientID:     "430d87d746fda65902940a414adeadfd",
		ClientSecret: os.Getenv("CLIENT_SECRET"),

		RedirectURL: localServer + "/v1/user/login/kakao",
		Endpoint: oauth2.Endpoint{
			AuthURL:  authKakao + "/oauth/authorize",
			TokenURL: authKakao + "/oauth/token",
		},
	}
)

type KakaoTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

func stateTokenGet() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func LoginKakaoAuthUrlGet(w http.ResponseWriter, r *http.Request) {
	reqId := getRequestId(w, r)
	logger.Debugf(reqId, "user/login/kakao/authurl GET started")

	authSession, _ := store.Get(r, "authSession")
	authSession.Options = &sessions.Options{
		Path:   "/v1/user/login/kakao",
		MaxAge: 300,
	}
	state := stateTokenGet()
	authSession.Values["state"] = state
	authSession.Save(r, w)

	logger.Debugf(reqId, "state created: %s", state)

	url := conf.AuthCodeURL(state, oauth2.AccessTypeOffline)

	// fmt.Println("clientID: ", os.Getenv("CLIENT_ID"))
	fmt.Println(url)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"oauthUrl": url})
}

func LoginKakaoGet(w http.ResponseWriter, r *http.Request) {
	reqId := getRequestId(w, r)
	logger.Debugf(reqId, "user/login/kakao GET started")

	authSession, _ := store.Get(r, "authSession")
	state := authSession.Values["state"]
	if state == nil {
		logger.Debugf(reqId, "state is %s", state)
		http.Redirect(w, r, "/v1/user/login/kakao/authurl", http.StatusFound)
		return
	}

	authSession.Options = &sessions.Options{
		MaxAge: -1,
	}
	authSession.Save(r, w)

	c := r.FormValue("code")
	fmt.Println(c)

	data := url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {os.Getenv("CLIENT_ID")},
		"redirect_uri":  {conf.RedirectURL},
		"code":          {c},
		"client_secret": {os.Getenv("CLIENT_SECRET")},
	}

	// 아래 코드 수정 필요
	fmt.Println(conf.Endpoint.TokenURL)
	resp, err := http.PostForm(conf.Endpoint.TokenURL, data)
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

	// Output the token information
	fmt.Printf("Access Token: %s\n", tokenResponse.AccessToken)
	fmt.Printf("Expires In: %d\n", tokenResponse.ExpiresIn)
	fmt.Printf("Refresh Token: %s\n", tokenResponse.RefreshToken)
	fmt.Printf("Scope: %s\n", tokenResponse.Scope)
	fmt.Printf("Token Type: %s\n", tokenResponse.TokenType)
}
