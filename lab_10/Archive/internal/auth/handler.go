package auth

import (
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"lab10/internal/config"
	"lab10/internal/email"
	"lab10/internal/users"
)

type Handler struct {
	userRepo *users.Repository
	emailSvc *email.Service
	cfg      *config.Config
}

func NewHandler(userRepo *users.Repository, emailSvc *email.Service, cfg *config.Config) *Handler {
	return &Handler{userRepo: userRepo, emailSvc: emailSvc, cfg: cfg}
}

type signUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type signInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// POST /api/auth/signup
func (h *Handler) SignUp(c *fiber.Ctx) error {
	var req signUpRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return fiber.NewError(fiber.StatusBadRequest, "email, password, and name are required")
	}
	if len(req.Password) < 8 {
		return fiber.NewError(fiber.StatusBadRequest, "password must be at least 8 characters")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return err
	}

	u, err := h.userRepo.Create(c.Context(), req.Email, string(hash), req.Name)
	if err != nil {
		// crude duplicate-email detection — pgx surfaces the constraint name
		if isPgUniqueViolation(err) {
			return fiber.NewError(fiber.StatusConflict, "email already taken")
		}
		return err
	}

	// fire-and-forget welcome email
	go func() {
		if err := h.emailSvc.SendTemplate(u.Email, "welcome.html", map[string]string{
			"Name":  u.Name,
			"Email": u.Email,
		}); err != nil {
			slog.Error("send welcome email", "err", err)
		}
	}()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":         u.ID,
		"email":      u.Email,
		"name":       u.Name,
		"created_at": u.CreatedAt,
	})
}

// POST /api/auth/signin
func (h *Handler) SignIn(c *fiber.Ctx) error {
	var req signInRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	u, err := h.userRepo.FindByEmail(c.Context(), req.Email)
	if errors.Is(err, users.ErrNotFound) {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
	}
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
	}

	ttl := time.Duration(h.cfg.JWTTTLHours) * time.Hour
	token, err := issueToken(h.cfg.JWTSecret, ttl, u.ID)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"token": token})
}

// GET /api/auth/me
func (h *Handler) Me(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	u, err := h.userRepo.FindByID(c.Context(), userID)
	if errors.Is(err, users.ErrNotFound) {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"id":         u.ID,
		"email":      u.Email,
		"name":       u.Name,
		"created_at": u.CreatedAt,
	})
}

func isPgUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	return strings.Contains(s, "23505") || strings.Contains(s, "unique")
}
