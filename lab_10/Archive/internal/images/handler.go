package images

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"lab10/internal/config"
)

var allowedMIME = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

type Handler struct {
	repo *Repository
	cfg  *config.Config
}

func NewHandler(repo *Repository, cfg *config.Config) *Handler {
	return &Handler{repo: repo, cfg: cfg}
}

// POST /api/images
func (h *Handler) Upload(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	file, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "field 'file' is required")
	}

	maxBytes := int64(h.cfg.MaxUploadSizeMB) * 1024 * 1024
	if file.Size > maxBytes {
		return fiber.NewError(fiber.StatusRequestEntityTooLarge,
			fmt.Sprintf("file exceeds %d MB limit", h.cfg.MaxUploadSizeMB))
	}

	mime := file.Header.Get("Content-Type")
	ext, allowed := allowedMIME[mime]
	if !allowed {
		return fiber.NewError(fiber.StatusBadRequest, "unsupported file type; allowed: jpeg, png, webp")
	}

	now := time.Now()
	dir := filepath.Join(h.cfg.UploadDir, now.Format("2006"), now.Format("01"))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create upload dir: %w", err)
	}

	storedName := uuid.New().String() + ext
	storedPath := filepath.Join(dir, storedName)

	if err := c.SaveFile(file, storedPath); err != nil {
		return fmt.Errorf("save file: %w", err)
	}

	img := &Image{
		OwnerID:      userID,
		OriginalName: file.Filename,
		StoredPath:   storedPath,
		MimeType:     mime,
		SizeBytes:    file.Size,
	}
	if err := h.repo.Create(c.Context(), img); err != nil {
		_ = os.Remove(storedPath)
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(img)
}

// GET /api/images/:id
func (h *Handler) Download(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid image id")
	}

	img, err := h.repo.FindByID(c.Context(), id)
	if errors.Is(err, ErrNotFound) {
		return fiber.NewError(fiber.StatusNotFound, "image not found")
	}
	if err != nil {
		return err
	}

	return c.SendFile(img.StoredPath)
}

// DELETE /api/images/:id
func (h *Handler) Delete(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid image id")
	}

	img, err := h.repo.FindByID(c.Context(), id)
	if errors.Is(err, ErrNotFound) {
		return fiber.NewError(fiber.StatusNotFound, "image not found")
	}
	if err != nil {
		return err
	}

	if img.OwnerID != userID {
		return fiber.NewError(fiber.StatusForbidden, "access denied")
	}

	if err := h.repo.Delete(c.Context(), id); err != nil {
		return err
	}

	_ = os.Remove(img.StoredPath)

	return c.SendStatus(fiber.StatusNoContent)
}
