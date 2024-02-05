package models

import (
	"log"
	"time"

	"gorm.io/gorm"
)

type CommonFields struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func SetDatabaseModel(db *gorm.DB) {
	err := db.AutoMigrate(&User{}, &Article{}, &Follow{}, &Favourite{})
	if err != nil {
		log.Fatal(err)
	}
}
