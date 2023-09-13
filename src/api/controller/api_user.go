package controller

import (
	"net/http"

	"github.com/jonggulee/go-login.git/src/logger"
)

func SignupUserPost(w http.ResponseWriter, r *http.Request) {
	reqId := getRequestId(w, r)
	logger.Debugf(reqId, "user/signup POST started")

}
