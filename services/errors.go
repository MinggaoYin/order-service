package services

import "errors"

var (
	ErrCannotCalculateDistance = errors.New("cannot calculate distance for given location")
	ErrOrderAlreadyTaken       = errors.New("order already taken")
)
