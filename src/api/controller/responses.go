package controller

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/jonggulee/como-backend/src/api/model"
	"github.com/jonggulee/como-backend/src/logger"
)

func addCommonHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
}

func addCORS(w http.ResponseWriter) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
}

func newOkResponse(w http.ResponseWriter, requestId, msg string) *model.Response {
	addCommonHeaders(w)
	addCORS(w)

	w.WriteHeader(http.StatusOK)

	return &model.Response{
		StatusCode:    200,
		StatusMessage: msg,
		RequestId:     requestId,
	}
}

func newResponse(w http.ResponseWriter, requestId string, code int, msg string) *model.Response {
	addCommonHeaders(w)
	addCORS(w)

	switch code {
	case 200:
		w.WriteHeader(http.StatusOK)
	case 400:
		w.WriteHeader(http.StatusBadRequest)
	case 403:
		w.WriteHeader(http.StatusForbidden)
	case 404:
		w.WriteHeader(http.StatusNotFound)
	case 500:
		w.WriteHeader(http.StatusInternalServerError)

	}

	return &model.Response{
		StatusCode:    code,
		StatusMessage: msg,
		RequestId:     requestId,
	}
}

func getRequestId(w http.ResponseWriter, r *http.Request) string {
	requestId := r.Header.Get("X-Request-Id")
	if requestId == "" {
		requestId = uuid.New().String()
	}
	w.Header().Set("X-Request-Id", requestId)
	return requestId
}

func writeResponse(reqId string, w http.ResponseWriter, response *model.Response) error {
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		logger.Errorf(reqId, "Failed to write response")
		return err
	}
	return nil
}
