package types

import (
	"hublish-be-go/internal/models"
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

type ShortUserQuery struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
	Followed bool   `json:"followed"`
}

type UserQuery struct {
	models.User
	Followed bool `json:"followed"`
}

type CommentQuery struct {
	ID              string `json:"id"`
	Body            string `json:"body"`
	CommentAuthorID string `json:"commentAuthor_id" gorm:"column:commentAuthor_id"`
	ArticleID       string `json:"article_id"`
	CommentAuthor   struct {
		Caid     string `json:"id"`
		Username string `json:"username"`
		Name     string `json:"name"`
		Image    string `json:"image"`
	} `json:"commentAuthor" gorm:"embedded"`
}
