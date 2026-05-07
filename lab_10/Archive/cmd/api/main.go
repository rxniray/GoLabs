package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	"lab10/internal/auth"
	"lab10/internal/config"
	"lab10/internal/database"
	"lab10/internal/documents"
	"lab10/internal/email"
	"lab10/internal/images"
	"lab10/internal/server"
	"lab10/internal/users"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func main() {
	initLogger()

	cfg, err := config.Load()
	if err != nil {
		slog.Error("load config", "err", err)
		os.Exit(1)
	}

	if err := runMigrations(cfg.DatabaseURL); err != nil {
		slog.Error("run migrations", "err", err)
		os.Exit(1)
	}

	ctx := context.Background()
	pool, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("connect to database", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	emailSvc, err := email.NewService(cfg)
	if err != nil {
		slog.Error("init email service", "err", err)
		os.Exit(1)
	}

	userRepo := users.NewRepository(pool)
	imageRepo := images.NewRepository(pool)

	authHandler := auth.NewHandler(userRepo, emailSvc, cfg)
	imagesHandler := images.NewHandler(imageRepo, cfg)
	docsHandler := documents.NewHandler()

	srv := server.New(cfg, authHandler, imagesHandler, docsHandler)

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.Start(); err != nil {
			slog.Error("server error", "err", err)
		}
	}()

	<-quit
	slog.Info("shutting down")
	if err := srv.Shutdown(); err != nil {
		slog.Error("shutdown error", "err", err)
	}
}

func runMigrations(dsn string) error {
	src, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("create migration source: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", src, "pgx5://"+stripScheme(dsn))
	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate up: %w", err)
	}
	slog.Info("migrations applied")
	return nil
}

// stripScheme removes an existing scheme prefix so we can inject "pgx5://"
func stripScheme(dsn string) string {
	for _, prefix := range []string{"postgres://", "postgresql://"} {
		if len(dsn) > len(prefix) && dsn[:len(prefix)] == prefix {
			return dsn[len(prefix):]
		}
	}
	return dsn
}

func initLogger() {
	env := os.Getenv("APP_ENV")
	var h slog.Handler
	if env == "production" {
		h = slog.NewJSONHandler(os.Stdout, nil)
	} else {
		h = slog.NewTextHandler(os.Stdout, nil)
	}
	slog.SetDefault(slog.New(h))
}
