package models

import (
	"github.com/lib/pq"
)

type Article struct {
	CommonFields   `gorm:"embeded"`
	Title          string         `json:"title"`
	Slug           string         `json:"slug" gorm:"unique;not null"`
	Content        string         `json:"content"`
	Tag            pq.StringArray `json:"tag" gorm:"type:text[]"`
	FavouriteCount uint           `json:"favouriteCount"`
	AuthorID       string         `json:"author_id"`
	Author         User           `gorm:"foreignKey:AuthorID"`
}
