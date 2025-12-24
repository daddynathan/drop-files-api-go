package file

import (
	"context"
	"database/sql"
)

type fileRepo struct {
	db *sql.DB
}

func NewFileRepo(db *sql.DB) FileRepo {
	return &fileRepo{db: db}
}

type FileRepo interface {
	Save(ctx context.Context, file *File) error
	FindByID(ctx context.Context, id string) (*File, error)
}

func (r *fileRepo) Save(ctx context.Context, file *File) error {
	query := `
		INSERT INTO files (id, filename, path, size, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx,
		query,
		file.ID,
		file.Filename,
		file.Path,
		file.Size,
		file.ExpiresAt,
		file.CreatedAt,
	)
	return err
}

func (r *fileRepo) FindByID(ctx context.Context, id string) (*File, error) {
	query := `
		SELECT id, filename, path, size, expires_at, created_at
		FROM files
		WHERE id = $1
	`
	file := &File{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&file.ID,
		&file.Filename,
		&file.Path,
		&file.Size,
		&file.ExpiresAt,
		&file.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return file, nil
}
