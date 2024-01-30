package types

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type ArticleQuery struct {
	ID        string         `json:"id"`
	Title     string         `json:"title"`
	Slug      string         `json:"slug"`
	Content   string         `json:"content"`
	Tags      pq.StringArray `json:"tags" gorm:"type:text[]"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
	AuthorID  string         `json:"author_id"`
	Author    struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Name     string `json:"name"`
		Bio      string `json:"bio"`
		Image    string `json:"image"`
	} `json:"author" gorm:"embedded"`
}
