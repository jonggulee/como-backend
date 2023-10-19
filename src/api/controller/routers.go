package controller

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, COMO\n\n"))
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		handler := route.HandlerFunc

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"AuthurlKakaoLoginUserGet",
		strings.ToUpper("GET"),
		"/v1/user/login/kakao/authurl",
		AuthurlKakaoLoginUserGet,
	},
	Route{
		"TokenKakaoLoginUserGet",
		strings.ToUpper("GET"),
		"/v1/user/login/kakao/token",
		TokenKakaoLoginUserGet,
	},
	Route{
		"DetailUserGet",
		strings.ToUpper("GET"),
		"/v1/user/detail",
		DetailUserGet,
	},
	Route{
		"DetailUserPost",
		strings.ToUpper("Post"),
		"/v1/user/detail",
		DetailUserPost,
	},
	Route{
		"WithdrawUserPost",
		strings.ToUpper("Post"),
		"/v1/user/withdraw",
		WithdrawUserPost,
	},
	Route{
		"ListEventGet",
		strings.ToUpper("GET"),
		"/v1/event",
		ListEventGet,
	},
}
