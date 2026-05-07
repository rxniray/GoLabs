package images

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("image not found")

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, img *Image) error {
	err := r.db.QueryRow(ctx,
		`INSERT INTO images (owner_id, original_name, stored_path, mime_type, size_bytes)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, created_at`,
		img.OwnerID, img.OriginalName, img.StoredPath, img.MimeType, img.SizeBytes,
	).Scan(&img.ID, &img.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert image: %w", err)
	}
	return nil
}

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*Image, error) {
	img := &Image{}
	err := r.db.QueryRow(ctx,
		`SELECT id, owner_id, original_name, stored_path, mime_type, size_bytes, created_at
		 FROM images WHERE id = $1`,
		id,
	).Scan(&img.ID, &img.OwnerID, &img.OriginalName, &img.StoredPath, &img.MimeType, &img.SizeBytes, &img.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find image: %w", err)
	}
	return img, nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM images WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete image: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
