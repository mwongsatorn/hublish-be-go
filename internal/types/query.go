package types

import (
	"hublish-be-go/internal/models"
	"time"

	"github.com/lib/pq"
)

type ArticleQuery struct {
	ID             string         `json:"id"`
	Title          string         `json:"title"`
	Slug           string         `json:"slug"`
	Content        string         `json:"content"`
	Tags           pq.StringArray `json:"tags" gorm:"type:text[]"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	AuthorID       string         `json:"author_id"`
	FavouriteCount int            `json:"favouriteCount"`
	Favourited     bool           `json:"favourited"`
	Author         struct {
		Aid      string `json:"id"`
		Username string `json:"username"`
		Name     string `json:"name"`
		Bio      string `json:"bio"`
		Image    string `json:"image"`
	} `json:"author" gorm:"embedded"`
}

type ShortUserQuery struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
	Followed bool   `json:"followed"`
}

type UserQuery struct {
	models.User
	Followed bool `json:"followed"`
}

type CommentQuery struct {
	ID              string    `json:"id"`
	Body            string    `json:"body"`
	CommentAuthorID string    `json:"commentAuthor_id" gorm:"column:commentAuthor_id"`
	ArticleID       string    `json:"article_id"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	CommentAuthor   struct {
		Caid     string `json:"id"`
		Username string `json:"username"`
		Name     string `json:"name"`
		Image    string `json:"image"`
	} `json:"commentAuthor" gorm:"embedded"`
}

type SearchQuery[T any] struct {
	TotalResults int `json:"total_results"`
	TotalPages   int `json:"total_pages"`
	Page         int `json:"page"`
	Results      []T `json:"results"`
}
