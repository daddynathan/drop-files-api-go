package file

import "errors"

var (
	ErrInvalidFilename = errors.New("filename cannot be empty or too long")
	ErrInvalidFileSize = errors.New("file size exceeds maximum allowed")
	ErrInvalidTTL      = errors.New("invalid TTL value")

	ErrFileNotFound = errors.New("file not found")
	ErrFileExpired  = errors.New("file expired")
)
