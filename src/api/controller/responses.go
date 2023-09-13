package controller

import (
	"net/http"

	"github.com/google/uuid"
)

func getRequestId(w http.ResponseWriter, r *http.Request) string {
	requestId := r.Header.Get("X-Request-Id")
	if requestId == "" {
		requestId = uuid.New().String()
	}
	w.Header().Set("X-Request-Id", requestId)
	return requestId
}
