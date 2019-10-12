package models

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrCannotUpdate = errors.New("cannot update due to conflict")
)
