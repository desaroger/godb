package common

import (
	"errors"
)

var (
	ErrDocumentAlreadyExists = errors.New("document already exists")
	ErrDocumentDoestNotExist = errors.New("document does not exist")
	ErrEmptyDocument         = errors.New("empty document")
	ErrInvalidId             = errors.New("invalid id")
)
