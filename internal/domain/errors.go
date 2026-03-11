package domain

import "errors"

// Sentinel domain errors — бизнес-ошибки, не привязанные к инфраструктуре.
var (
	ErrNotFound           = errors.New("resource not found")
	ErrAlreadyExists      = errors.New("resource already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrForbidden          = errors.New("access denied")
	ErrTokenInvalid       = errors.New("token invalid")
)
