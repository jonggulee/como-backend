package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/jonggulee/como-backend/src/api/model"
	"github.com/jonggulee/como-backend/src/config"
	"github.com/jonggulee/como-backend/src/constants"
	dbController "github.com/jonggulee/como-backend/src/db/controller"
	"github.com/jonggulee/como-backend/src/logger"
)

func DecodeJwt(reqId, auth string) (*model.Session, error) {
	logger.Debugf(reqId, "Try to decode JWT token")

	claims := jwt.MapClaims{}

	if strings.HasPrefix(strings.ToLower(auth), "bearer ") {
		auth = auth[7:]
	}

	keyFunc := GetKeyFunc()

	token, err := jwt.ParseWithClaims(auth, claims, keyFunc)
	if err != nil {
		logger.Errorf(reqId, "Failed to parse JWT token: %s, due to %s", auth, err)
		return nil, err
	}

	if !token.Valid {
		logger.Errorf(reqId, "Invalid JWT token: %s", auth)
		return nil, err
	}

	sessionId, ok := claims["session"].(string)
	if !ok {
		logger.Errorf(reqId, "Failed to get session_id from JWT token: %s", auth)
		return nil, err
	}

	sessionInfo, err := dbController.SessionSelect(config.AppCtx.Db.Db, reqId, sessionId)
	if err != nil {
		logger.Errorf(reqId, "Failed to select session by session: %s, due to %s", sessionId, err)
		return nil, err
	}

	if sessionInfo.Id == 0 {
		logger.Errorf(reqId, "Session not found: %s", sessionId)
		return nil, errors.New("session not found")
	}

	return sessionInfo, nil
}

func GetKeyFunc() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(constants.JWTPSK), nil
	}
}

func GetIntegerFromPathParameters(pathParameters map[string]string, key string) (int, error) {
	if value, ok := pathParameters[key]; ok {
		i, err := strconv.Atoi(value)
		if err != nil {
			return 0, err
		}
		return i, nil
	}
	return 0, fmt.Errorf("failed to found %s from path parameters", key)
}
