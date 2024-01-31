package models

type Favourite struct {
	CommonFields
	UserID    string `json:"user_id"`
	ArticleID string `json:"article_id"`
	User      User
	Article   Article
}
