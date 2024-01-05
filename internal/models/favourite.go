package models

type Favourite struct {
	CommonFields `gorm:"embeded"`
	UserID       string `json:"user_id"`
	ArticleID    string `json:"article_id"`
	User         User
	Article      Article
}
