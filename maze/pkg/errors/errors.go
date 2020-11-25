package errors

import "errors"

var (
	ErrResponseEncoding     = errors.New("an error occurred encoding response")
	ErrMissingBodyContent   = errors.New("missing content")
	ErrMalformedBodyContent = errors.New("malformed content")
)
