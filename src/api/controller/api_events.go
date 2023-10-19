package controller

import (
	"net/http"

	"github.com/jonggulee/go-login.git/src/config"
	"github.com/jonggulee/go-login.git/src/constants"
	dbController "github.com/jonggulee/go-login.git/src/db/controller"
	"github.com/jonggulee/go-login.git/src/logger"
)

func ListEventGet(w http.ResponseWriter, r *http.Request) {
	reqId := getRequestId(w, r)
	logger.Debugf(reqId, "event/list GET started")

	events, err := dbController.EventListSelect(config.AppCtx.Db.Db, reqId)
	if err != nil {
		logger.Errorf(reqId, "Failed to select from event... %s", err)
		resp := newResponse(w, reqId, 500, "Internal Server Error")
		writeResponse(reqId, w, resp)
		return
	}

	resp := newOkResponse(w, reqId, constants.BASICOK)
	resp.Events = events
	writeResponse(reqId, w, resp)

	// writeResponse(w, r, ht
}
