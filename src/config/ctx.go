package config

import (
	dbModel "github.com/jonggulee/go-login.git/src/db/model"
)

type Ctx struct {
	Cfg *Config
	Db  *dbModel.DbConnection
}

var (
	AppCtx = &Ctx{
		Cfg: &Config{},
		Db:  &dbModel.DbConnection{},
	}
)
