package errs

import "errors"

var (
	ErrNoSeats          = errors.New("no seats available")
	ErrNotFound         = errors.New("not found")
	ErrAlreadyConfirmed = errors.New("already confirmed")
	ErrCanceled         = errors.New("canceled")
	ErrExpired          = errors.New("expired")
	ErrInvalid          = errors.New("invalid input")
)
