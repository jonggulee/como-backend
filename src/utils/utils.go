package utils

import (
	"errors"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/jonggulee/go-login.git/src/constants"
	"github.com/jonggulee/go-login.git/src/logger"
)

func DecodeJwt(reqId, auth string) (*model.SessionInfo, error) {
	logger.Debugf(reqId, "Start decoding JWT token: %s", auth)

	claims := jwt.MapClaims{}

	auth = strings.Replace(auth, "Bearer ", "", 1)
	auth = strings.Replace(auth, "bearer ", "", 1)

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

}

func GetKeyFunc() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(constants.JWTPSK), nil
	}
}
