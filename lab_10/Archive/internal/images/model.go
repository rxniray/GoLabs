package images

import (
	"time"

	"github.com/google/uuid"
)

type Image struct {
	ID           uuid.UUID
	OwnerID      uuid.UUID
	OriginalName string
	StoredPath   string
	MimeType     string
	SizeBytes    int64
	CreatedAt    time.Time
}
