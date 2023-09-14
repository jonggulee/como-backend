package controller

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/jonggulee/go-login.git/src/logger"
	"golang.org/x/oauth2"
)

const (
// localServer = "http://localhost:8080"
// authKakao   = "https://kauth.kakao.com"
)

var (
	store = sessions.NewCookieStore([]byte("secret"))

	conf = &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		// ClientID: "430d87d746fda65902940a414adeadfd",
		// ClientID:     os.Getenv("ClientID"),

		// RedirectURL:  localServer + "/v1/user/signup",
		RedirectURL: "http://localhost:8080/v1/user/signup",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://kauth.kakao.com/oauth/authorize",
			TokenURL: "https://kauth.kakao.com/oauth/token",
			// AuthURL:  authKakao + "/oauth/authorize",
			// TokenURL: authKakao + "/oauth/token",
		},
	}
)

func stateTokenGet() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

// func sessionStore(w http.ResponseWriter, r *http.Request, state string) {
// 	session, _ := store.Get(r, "session")
// 	session.Options = &sessions.Options{
// 		Path:   "/v1/user/login/kakao/authurl",
// 		MaxAge: 300,
// 	}
// 	session.Values["state"] = state
// 	session.Save(r, w)
// }

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

	// fmt.Println("clientID: ", os.Getenv("CLIENT_ID"), "clientSecret", conf.ClientSecret)
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
		return
	}

	// authSession.Options = &sessions.Options{
	// 	MaxAge: -1,
	// }
	// authSession.Save(r, w)

	fmt.Println("code: ", r.Form.Get("code"))
}

// func authorizeCodeGet(w http.ResponseWriter, r *http.Request) {
// 	reqId := getRequestId(w, r)
// 	logger.Debugf(reqId, "user/signup/authorizeCode GET started")

// 	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
// 	fmt.Println(url)
// 	http.Redirect(w, r, url, http.StatusFound)
// }

// func SingupTokenPost(w http.ResponseWriter, r *http.Request) {
// 	reqId := getRequestId(w, r)
// 	logger.Debugf(reqId, "user/signup/token POST started")

// 	s := r.FormValue("state")
// 	if s != "state" {
// 		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", "state", s)
// 		http.Redirect(w, r, "/", http.StatusFound)
// 		return
// 	}
// 	logger.Debugf(reqId, "state: %s", s)
// 	c := r.FormValue("code")
// 	logger.Debugf(reqId, "code: %s", c)

// }

// func SignupUserGet(w http.ResponseWriter, r *http.Request) {
// 	reqId := getRequestId(w, r)
// 	logger.Debugf(reqId, "user/signup GET started")

// 	authorizeCodeGet(w, r)

// }
