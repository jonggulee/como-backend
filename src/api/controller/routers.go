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
		strings.ToUpper("Get"),
		"/v1/user/login/kakao/authurl",
		AuthurlKakaoLoginUserGet,
	},
	Route{
		"TokenKakaoLoginUserGet",
		strings.ToUpper("Get"),
		"/v1/user/login/kakao/token",
		TokenKakaoLoginUserGet,
	},
	Route{
		"DetailUserGet",
		strings.ToUpper("Get"),
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
		"WithdrawUserDelete",
		strings.ToUpper("Delete"),
		"/v1/user/withdraw",
		WithdrawUserDelete,
	},
	Route{
		"EventGet",
		strings.ToUpper("Get"),
		"/v1/event",
		EventGet,
	},
	Route{
		"EventCreatePost",
		strings.ToUpper("Post"),
		"/v1/event",
		EventCreatePost,
	},
	Route{
		"EventEditPost",
		strings.ToUpper("Post"),
		"/v1/event/{eventId}",
		EventEditPost,
	},
	Route{
		"EventDelete",
		strings.ToUpper("Delete"),
		"/v1/event/{eventId}",
		EventDelete,
	},
}
