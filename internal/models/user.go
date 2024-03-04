package models

import "github.com/lib/pq"

type User struct {
	CommonFields
	Username       string         `json:"username" gorm:"not null;unique"`
	Password       string         `json:"-"`
	Email          string         `json:"email" gorm:"not null;unique"`
	Name           string         `json:"name" gorm:"default:''"`
	Bio            string         `json:"bio" gorm:"default:''"`
	Image          string         `json:"image" gorm:"default:'https://api.dicebear.com/7.x/thumbs/svg?radius=50&size=240&backgroundColor=f0f0f0&mouth=variant2&shapeColor=16a34a'"`
	RefreshTokens  pq.StringArray `json:"-" gorm:"type:text[];default:'{}'"`
	FollowerCount  uint           `json:"followerCount" gorm:"default:0"`
	FollowingCount uint           `json:"followingCount" gorm:"default:0"`
}
