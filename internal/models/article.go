package models

import (
	"github.com/lib/pq"
)

type Article struct {
	CommonFields   `gorm:"embeded"`
	Title          string         `json:"title"`
	Slug           string         `json:"slug" gorm:"unique;not null"`
	Content        string         `json:"content"`
	Tags           pq.StringArray `json:"tags" gorm:"type:text[]"`
	FavouriteCount uint           `json:"favouriteCount"`
	AuthorID       string         `json:"author_id"`
	Author         User           `json:"author" gorm:"foreignKey:AuthorID"`
}
