package config

import (
	dbModel "github.com/jonggulee/como-backend/src/db/model"
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
