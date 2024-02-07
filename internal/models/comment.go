package models

type Comment struct {
	CommonFields
	Body            string  `json:"body"`
	CommentAuthorID string  `json:"commentAuthor_id" gorm:"column:commentAuthor_id"`
	ArticleID       string  `json:"article_id"`
	User            User    `json:"user" gorm:"foreignKey:CommentAuthorID;"`
	Article         Article `json:"article"`
}
