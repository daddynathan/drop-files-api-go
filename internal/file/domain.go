package file

import (
	"time"
)

type File struct {
	ID        string
	Filename  string
	Path      string
	Size      int64
	ExpiresAt time.Time
	CreatedAt time.Time
}

func NewFile(id, originalName, storagePath string, size int64, expiresAt time.Time) *File {
	return &File{
		ID:        id,
		Filename:  originalName,
		Path:      storagePath,
		Size:      size,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}
}
