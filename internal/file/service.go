package file

import (
	"DropFiles/config"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/google/uuid"
)

type FileService struct {
	repo    FileRepo
	storage Storage
	config  config.FileConfig
}

func NewFileService(repo FileRepo, storage Storage, cfg config.FileConfig) *FileService {
	return &FileService{
		repo:    repo,
		storage: storage,
		config:  cfg,
	}
}

func (s *FileService) Upload(ctx context.Context, filename string, r io.Reader, size int64, ttlHours int) (string, error) {
	if filename == "" || len(filename) < s.config.MinFilenameLen || len(filename) > s.config.MaxFilenameLen {
		return "", ErrInvalidFilename
	}
	if size <= 0 || size > s.config.MaxFileSizeBytes {
		return "", ErrInvalidFileSize
	}
	defaultTTLHours := s.config.DefaultTTLHours
	if ttlHours > 0 && ttlHours <= s.config.MaxTTLHours {
		defaultTTLHours = ttlHours
	}
	id := uuid.New().String()
	storagePath := fmt.Sprintf("file_%s_%d", id, size)
	expiresAt := time.Now().Add(time.Duration(defaultTTLHours) * time.Hour)
	file := NewFile(id, filename, storagePath, size, expiresAt)
	if err := s.storage.Save(file.Path, r); err != nil {
		return "", fmt.Errorf("failed to save file to disk: %w", err)
	}
	if err := s.repo.Save(ctx, file); err != nil {
		os.Remove(s.storage.FullPath(file.Path))
		return "", fmt.Errorf("failed to save to db: %w", err)
	}
	return file.ID, nil
}

func (s *FileService) Download(ctx context.Context, id string) (string, error) {
	if id == "" {
		return "", ErrFileNotFound
	}
	f, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrFileNotFound
		}
		return "", err
	}
	if time.Now().After(f.ExpiresAt) {
		// Опционально: удалять просроченные файлы
		return "", ErrFileExpired
	}
	return f.Path, nil
}
