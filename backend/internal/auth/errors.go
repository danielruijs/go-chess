package auth

import "errors"

var ErrInvalidCredentials = errors.New("invalid username or password")

type UserRegistrationError struct {
	Message string
}

func (e *UserRegistrationError) Error() string {
	return e.Message
}

func NewUserRegistrationError(msg string) error {
	return &UserRegistrationError{Message: msg}
}
