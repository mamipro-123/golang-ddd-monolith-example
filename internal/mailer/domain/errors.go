package domain

import "errors"

// ErrInvalidEmail represents an error when an email is not valid.
var ErrInvalidEmail = errors.New("invalid email address")
