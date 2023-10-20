package controller

import (
	"net/http"

	"github.com/jonggulee/go-login.git/src/config"
	"github.com/jonggulee/go-login.git/src/constants"
	dbController "github.com/jonggulee/go-login.git/src/db/controller"
	"github.com/jonggulee/go-login.git/src/logger"
)

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

	// if len(*events) == 0 {
	// 	logger.Errorf(reqId, "No event found")
	// 	resp := newResponse(w, reqId, 404, "No event found")
	// 	writeResponse(reqId, w, resp)
	// 	return
	// }

	resp := newOkResponse(w, reqId, constants.BASICOK)
	resp.Events = events
	writeResponse(reqId, w, resp)
}

func EventPost(w http.ResponseWriter, r *http.Request) {
	reqId := getRequestId(w, r)
	logger.Debugf(reqId, "event POST started")

	// eventReq := &model.EventRequest{}

	resp := newOkResponse(w, reqId, constants.BASICOK)
	writeResponse(reqId, w, resp)
}
