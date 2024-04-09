package errors

import "errors"

var (
	ErrUnauthorized   = errors.New("unauthorized access")
	ErrForbidden      = errors.New("forbidden access")
	ErrBannerNotFound = errors.New("banner not found")
)
