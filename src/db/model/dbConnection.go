package model

import "gorm.io/gorm"

type DbConnection struct {
	Db *gorm.DB
}
