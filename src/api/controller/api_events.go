package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jonggulee/como-backend/src/api/model"
	"github.com/jonggulee/como-backend/src/config"
	"github.com/jonggulee/como-backend/src/constants"
	dbController "github.com/jonggulee/como-backend/src/db/controller"
	"github.com/jonggulee/como-backend/src/logger"
	"github.com/jonggulee/como-backend/src/utils"
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
		return fmt.Errorf("user is not admin")
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

func EventCreatePost(w http.ResponseWriter, r *http.Request) {
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

func EventEditPost(w http.ResponseWriter, r *http.Request) {
	reqId := getRequestId(w, r)
	logger.Debugf(reqId, "event/edit/{eventId} POST started")

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

	pathParameters := mux.Vars(r)
	if pathParameters == nil {
		logger.Errorf(reqId, "Failed to get path parameters")
		resp := newResponse(w, reqId, 400, "Bad Request")
		writeResponse(reqId, w, resp)
		return
	}

	eventReq := &model.EventCreateRequest{}
	eventReq.UpdateUserId = decodedJwt.UserId

	err = json.NewDecoder(r.Body).Decode(eventReq)
	if err != nil {
		logger.Errorf(reqId, "Failed to decode request body %s", err)
		resp := newResponse(w, reqId, 400, "Bad Request")
		writeResponse(reqId, w, resp)
		return
	}

	eventId, err := utils.GetIntegerFromPathParameters(pathParameters, "eventId")
	if err != nil {
		logger.Errorf(reqId, "Failed to get eventId from path parameters, due to %s", err)
		resp := newResponse(w, reqId, 400, "Bad Request")
		writeResponse(reqId, w, resp)
		return
	}

	err = dbController.EventSelectById(config.AppCtx.Db.Db, reqId, eventId)
	if err != nil {
		logger.Errorf(reqId, "Failed to select event by id %d, due to %s", eventId, err)
		resp := newResponse(w, reqId, 400, "Bad Request")
		writeResponse(reqId, w, resp)
		return
	}

	err = dbController.EventEditPost(config.AppCtx.Db.Db, reqId, eventId, eventReq)
	if err != nil {
		logger.Errorf(reqId, "Failed to update event... %d, due to %s", eventId, err)
		resp := newResponse(w, reqId, 500, "Internal Server Error")
		writeResponse(reqId, w, resp)
		return
	}

	resp := newOkResponse(w, reqId, constants.BASICOK)
	writeResponse(reqId, w, resp)
}

func EventDelete(w http.ResponseWriter, r *http.Request) {
	reqId := getRequestId(w, r)
	logger.Debugf(reqId, "event/{eventId} DELETE started")

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

	pathParameters := mux.Vars(r)
	if pathParameters == nil {
		logger.Errorf(reqId, "Failed to get path parameters")
		resp := newResponse(w, reqId, 400, "Bad Request")
		writeResponse(reqId, w, resp)
		return
	}

	eventId, err := utils.GetIntegerFromPathParameters(pathParameters, "eventId")
	if err != nil {
		logger.Errorf(reqId, "Failed to get eventId from path parameters, due to %s", err)
		resp := newResponse(w, reqId, 400, "Bad Request")
		writeResponse(reqId, w, resp)
		return
	}

	err = dbController.EventSelectById(config.AppCtx.Db.Db, reqId, eventId)
	if err != nil {
		logger.Errorf(reqId, "Failed to select event by id %d, due to %s", eventId, err)
		resp := newResponse(w, reqId, 400, "Bad Request")
		writeResponse(reqId, w, resp)
		return
	}

	err = dbController.EventDelete(config.AppCtx.Db.Db, reqId, eventId, decodedJwt.UserId)
	if err != nil {
		logger.Errorf(reqId, "Failed to delete event... %d, due to %s", eventId, err)
		resp := newResponse(w, reqId, 500, "Internal Server Error")
		writeResponse(reqId, w, resp)
		return
	}

	resp := newOkResponse(w, reqId, constants.BASICOK)
	writeResponse(reqId, w, resp)
}
