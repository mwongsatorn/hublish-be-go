package models

type User struct {
	CommonFields   `gorm:"embeded"`
	Username       string `json:"username" gorm:"not null;unique"`
	Password       string `json:"password"`
	Email          string `json:"email" gorm:"not null;unique"`
	Name           string `json:"name" gorm:"default:''"`
	Bio            string `json:"bio" gorm:"default:''"`
	Image          string `json:"image" gorm:"default:''"`
	FollowerCount  uint   `json:"followerCount" gorm:"default:0"`
	FollowingCount uint   `json:"followingCount" gorm:"default:0"`
}
