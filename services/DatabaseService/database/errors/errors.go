package errors

import "errors"

var (
	ErrUnsupported	= errors.New("unsupported")
	ErrNoDocuments	= errors.New("no documents found")
)

