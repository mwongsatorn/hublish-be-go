package models

type Comment struct {
	CommonFields
	Body            string `json:"body"`
	CommentAuthorID string `json:"commentAuthor_id" gorm:"column:commentAuthor_id"`
	ArticleID       string `json:"article_id"`
	User            User   `gorm:"foreignKey:CommentAuthorID;"`
	Article         Article
}
