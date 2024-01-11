package validator

import "github.com/go-playground/validator/v10"

type SignUpRequest struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
	Email    string `validate:"required,email"`
}

var V = validator.New(validator.WithRequiredStructEnabled())
