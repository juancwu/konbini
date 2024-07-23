package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/juancwu/konbini/util"
)

// customValidator represents the custom validator the echo server uses to validate data.
type customValidator struct {
	validator *validator.Validate
}

// Validate method to satisfy the echo validator interface.
func (cv *customValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// validatePassword is a custom validator to validate passwords.
// It checks that a password is at least 12 characters long,
// it contains at least one special character,
// at least one uppercase letter,
// at least one lowercase letter,
// and at least one digit.
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	return util.ValidatePassword(password)
}
