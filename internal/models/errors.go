package models

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrAlreadyExists   = errors.New("already exists")
	ErrInvalidInput    = errors.New("invalid input")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrForbidden       = errors.New("forbidden")
	ErrTenantRequired  = errors.New("tenant required")
	ErrTenantMismatch  = errors.New("tenant mismatch")
	ErrValidation      = errors.New("validation failed")
	ErrConflict        = errors.New("conflict")
	ErrDatabase        = errors.New("database error")
	ErrEncoding        = errors.New("encoding error")
	ErrDecoding        = errors.New("decoding error")
	ErrUnsupportedType = errors.New("unsupported type")
)

func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

func IsAlreadyExists(err error) bool {
	return errors.Is(err, ErrAlreadyExists)
}

func IsInvalidInput(err error) bool {
	return errors.Is(err, ErrInvalidInput)
}

func IsUnauthorized(err error) bool {
	return errors.Is(err, ErrUnauthorized)
}

func IsForbidden(err error) bool {
	return errors.Is(err, ErrForbidden)
}

func IsTenantRequired(err error) bool {
	return errors.Is(err, ErrTenantRequired)
}

func IsValidation(err error) bool {
	return errors.Is(err, ErrValidation)
}

func IsDatabase(err error) bool {
	return errors.Is(err, ErrDatabase)
}
