package models

import (
	"log"
	"time"

	"gorm.io/gorm"
)

type CommonFields struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

func SetDatabaseModel(db *gorm.DB) {
	err := db.AutoMigrate(&User{}, &Article{}, &Follow{}, &Favourite{}, &Comment{})
	if err != nil {
		log.Fatal(err)
	}
}
