package controller

import (
	"encoding/json"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/jonggulee/go-login.git/src/api/model"
	"github.com/jonggulee/go-login.git/src/config"
	"github.com/jonggulee/go-login.git/src/constants"
	dbController "github.com/jonggulee/go-login.git/src/db/controller"
	"github.com/jonggulee/go-login.git/src/logger"
	"github.com/jonggulee/go-login.git/src/utils"
)

func checkAdminPermission(reqId string, userId int) error {
	logger.Debugf(reqId, "Try to check permission")

	user, err := dbController.UserDetailSelect(config.AppCtx.Db.Db, reqId, userId)
	if err != nil {
		logger.Errorf(reqId, "Failed to select from user... %s", err)
		return err
	}

	if user.Role != constants.ADMIN {
		logger.Errorf(reqId, "User is not admin")
		return err
	}

	return nil
}

func EventGet(w http.ResponseWriter, r *http.Request) {
	reqId := getRequestId(w, r)
	logger.Debugf(reqId, "event/list GET started")

	events, err := dbController.EventSelect(config.AppCtx.Db.Db, reqId)
	if err != nil {
		logger.Errorf(reqId, "Failed to select from event... %s", err)
		resp := newResponse(w, reqId, 500, "Internal Server Error")
		writeResponse(reqId, w, resp)
		return
	}

	resp := newOkResponse(w, reqId, constants.BASICOK)
	resp.Events = events
	writeResponse(reqId, w, resp)
}

func EventPost(w http.ResponseWriter, r *http.Request) {
	reqId := getRequestId(w, r)
	logger.Debugf(reqId, "event POST started")

	token := r.Header.Get("Authorization")
	if token == "" {
		logger.Errorf(reqId, "Failed to get Authorization header")
		resp := newResponse(w, reqId, 400, "Bad Request")
		writeResponse(reqId, w, resp)
		return
	}

	decodedJwt, err := utils.DecodeJwt(reqId, token)
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			// 토큰 만료 에러 처리
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

	err = checkAdminPermission(reqId, decodedJwt.UserId)
	if err != nil {
		logger.Errorf(reqId, "Failed to check admin permission %s", err)
		resp := newResponse(w, reqId, 403, "Forbidden")
		writeResponse(reqId, w, resp)
		return
	}

	eventReq := &model.EventCreateRequest{}
	eventReq.CreateUserId = decodedJwt.UserId
	eventReq.UpdateUserId = decodedJwt.UserId

	err = json.NewDecoder(r.Body).Decode(eventReq)
	if err != nil {
		logger.Errorf(reqId, "Failed to decode request body %s", err)
		resp := newResponse(w, reqId, 400, "Bad Request")
		writeResponse(reqId, w, resp)
		return
	}

	err = dbController.EventCreateInsert(config.AppCtx.Db.Db, reqId, eventReq)
	if err != nil {
		logger.Errorf(reqId, "Failed to insert into event... %v, due to %s", eventReq, err)
		resp := newResponse(w, reqId, 500, "Internal Server Error")
		writeResponse(reqId, w, resp)
		return
	}

	resp := newOkResponse(w, reqId, constants.BASICOK)
	writeResponse(reqId, w, resp)
}
