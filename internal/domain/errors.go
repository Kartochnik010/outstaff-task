package domain

import "errors"

var (
	ErrInvalidDate = errors.New("invalid date")
	ErrInternal    = errors.New("internal error")

	ErrMusicNotFound = errors.New("music not found")
)
