package storage

import "errors"

var (
	ErrNoID     = errors.New("No ID specified")
	ErrNotFound = errors.New("Record not found")
)
