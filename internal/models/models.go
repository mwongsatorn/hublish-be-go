package models

import (
	"gorm.io/gorm"
)

var db *gorm.DB

func SetDatabaseModel(gormDb *gorm.DB) {
	db = gormDb
}
