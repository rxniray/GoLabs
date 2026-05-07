package server

import (
	"errors"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"lab10/internal/auth"
	"lab10/internal/config"
	"lab10/internal/documents"
	"lab10/internal/images"
)

type Server struct {
	app *fiber.App
	cfg *config.Config
}

func New(
	cfg *config.Config,
	authHandler *auth.Handler,
	imagesHandler *images.Handler,
	docsHandler *documents.Handler,
) *Server {
	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler,
	})

	app.Use(recover.New())
	app.Use(requestLogger())

	api := app.Group("/api")

	// Auth routes
	a := api.Group("/auth")
	a.Post("/signup", authHandler.SignUp)
	a.Post("/signin", authHandler.SignIn)
	a.Get("/me", auth.JWTMiddleware(cfg.JWTSecret), authHandler.Me)

	// Images routes (all protected)
	img := api.Group("/images", auth.JWTMiddleware(cfg.JWTSecret))
	img.Post("", imagesHandler.Upload)
	img.Get("/:id", imagesHandler.Download)
	img.Delete("/:id", imagesHandler.Delete)

	// Documents routes (protected)
	doc := api.Group("/documents", auth.JWTMiddleware(cfg.JWTSecret))
	doc.Post("/generate", docsHandler.Generate)

	return &Server{app: app, cfg: cfg}
}

func (s *Server) Start() error {
	addr := ":" + s.cfg.AppPort
	slog.Info("starting server", "addr", addr, "env", s.cfg.AppEnv)
	return s.app.Listen(addr)
}

func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}

func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	msg := "internal server error"

	var fe *fiber.Error
	if errors.As(err, &fe) {
		code = fe.Code
		msg = fe.Message
	}

	errCode := httpCodeToErrCode(code)
	return c.Status(code).JSON(fiber.Map{
		"error": fiber.Map{
			"code":    errCode,
			"message": msg,
		},
	})
}

func httpCodeToErrCode(code int) string {
	m := map[int]string{
		400: "BAD_REQUEST",
		401: "UNAUTHORIZED",
		403: "FORBIDDEN",
		404: "NOT_FOUND",
		409: "CONFLICT",
		413: "PAYLOAD_TOO_LARGE",
		500: "INTERNAL_ERROR",
	}
	if s, ok := m[code]; ok {
		return s
	}
	return "ERROR"
}

func requestLogger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${time} | ${method} ${path} | ${status} | ${latency}\n",
		TimeFormat: time.RFC3339,
	})
}
