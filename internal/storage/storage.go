package storage

import "errors"

var (
    ErrCarNotFound = errors.New("car not found")
    ErrCarExists   = errors.New("car exists")
	ErrPersonNotFound = errors.New("person not found")
    ErrPersonExists   = errors.New("person exists")
)