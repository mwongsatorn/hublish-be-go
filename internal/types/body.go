package types

import "github.com/lib/pq"

type LoginRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignUpRequestBody struct {
	Username string `validate:"required"`
	Password string `validate:"required" validate:"min=8"`
	Email    string `validate:"required,email"`
}

type ChangeProfileRequestBody struct {
	Name  *string `json:"name" validate:"max=70"`
	Bio   *string `json:"bio" validate:"max=160"`
	Image *string `json:"image"`
}

type ChangePasswordRequestBody struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword" validate:"min=8"`
}

type ChangeEmailRequestBody struct {
	NewEmail string `json:"newEmail" validate:"email"`
	Password string `json:"password" `
}

type CreateArticleRequestBody struct {
	Title   string         `json:"title" validate:"required,max=70"`
	Content string         `json:"content" validate:"required,max=1500"`
	Tags    pq.StringArray `json:"tags" validate:"omitempty,dive,max=20" gorm:"type:text[];default:'{}'"`
}

type EditArticleRequestBody struct {
	Title   *string        `json:"title" validate:"omitempty,max=70"`
	Content *string        `json:"content" validate:"omitempty,max=1500"`
	Tags    pq.StringArray `json:"tags" validate:"omitempty,dive,max=20" gorm:"type:text[]"`
	Slug    string
}

type AddCommentRequestBody struct {
	Body string `json:"body" validate:"required,max=500"`
}
