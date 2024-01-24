package models

import "github.com/lib/pq"

type User struct {
	CommonFields   `gorm:"embeded"`
	Username       string         `json:"username" gorm:"not null;unique"`
	Password       string         `json:"-"`
	Email          string         `json:"email" gorm:"not null;unique"`
	Name           string         `json:"name" gorm:"default:''"`
	Bio            string         `json:"bio" gorm:"default:''"`
	Image          string         `json:"image" gorm:"default:''"`
	RefreshTokens  pq.StringArray `json:"-" gorm:"type:text[];default:'{}'"`
	FollowerCount  uint           `json:"followerCount" gorm:"default:0"`
	FollowingCount uint           `json:"followingCount" gorm:"default:0"`
}
