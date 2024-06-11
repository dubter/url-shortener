package domain

import "errors"

var (
	ErrURLNotFound = errors.New("requested resource is not found")
)
