package models

type Follow struct {
	CommonFields `gorm:"embeded"`
	FollowingID  string `json:"following_id"`
	FollowerID   string `json:"follower_id"`
	Following    User   `gorm:"foreignKey:FollowingID"`
	Follower     User   `gorm:"foreignKey:FollowerID"`
}
