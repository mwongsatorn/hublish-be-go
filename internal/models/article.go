package models

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Article struct {
	CommonFields   `gorm:"embeded"`
	Title          string                    `json:"title"`
	Slug           string                    `json:"slug" gorm:"unique;not null"`
	Content        string                    `json:"content"`
	Tag            pgtype.Array[pgtype.Text] `json:"tag" gorm:"type:text[]"`
	FavouriteCount uint                      `json:"favouriteCount"`
	AuthorID       uint                      `json:"author_id"`
	Author         User                      `gorm:"foreignKey:AuthorID"`
}
