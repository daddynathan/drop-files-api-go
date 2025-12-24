package file

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Storage interface {
	Save(path string, r io.Reader) error
	FullPath(path string) string
}

type LocalStorage struct {
	uploadDir string
}

func NewLocalStorage(uploadDir string) Storage {
	abs, _ := filepath.Abs(uploadDir)
	return &LocalStorage{uploadDir: abs}
}

func (s *LocalStorage) Save(path string, r io.Reader) error {
	cleanPath := filepath.Clean(path)
	if strings.Contains(cleanPath, "..") {
		return os.ErrInvalid
	}
	fullPath := filepath.Join(s.uploadDir, cleanPath)
	if !strings.HasPrefix(fullPath, s.uploadDir) {
		return os.ErrInvalid
	}
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}
	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, r)
	return err
}

func (s *LocalStorage) FullPath(path string) string {
	return filepath.Join(s.uploadDir, filepath.Clean(path))
}
