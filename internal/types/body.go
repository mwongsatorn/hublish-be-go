package types

type LoginRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignUpRequestBody struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
	Email    string `validate:"required,email"`
}

type UpdateProfileRequestBody struct {
	Name  *string `json:"name" validate:"max=70"`
	Bio   *string `json:"bio" validate:"max=160"`
	Image *string `json:"image"`
}

type ChangePasswordRequestBody struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}
