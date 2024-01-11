package models

import (
	"log"
	"time"

	"gorm.io/gorm"
)

type CommonFields struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func SetDatabaseModel(db *gorm.DB) {
	err := db.AutoMigrate(&User{}, &Article{}, &Follow{}, &Favourite{})
	if err != nil {
		log.Fatal(err)
	}
}
