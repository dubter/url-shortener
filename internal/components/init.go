package components

import (
	"url-shortner/internal/ports"
	"url-shortner/internal/services/encoder"
	"url-shortner/internal/services/render"
	"url-shortner/internal/services/url_shortener"

	"log/slog"
	"os"
	"url-shortner/internal/config"
	"url-shortner/internal/storage/pg"
	"url-shortner/internal/storage/redis"
	"url-shortner/pkg/logger/slogpretty"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type Components struct {
	HttpServer *ports.Server
	Postgres   *pg.Postgres
	Redis      *redis.Redis
}

func InitComponents(cfg *config.Config, logger *slog.Logger) (*Components, error) {
	postgres, err := pg.New(cfg.Postgres.PostgresURL)
	if err != nil {
		return nil, err
	}

	rds, err := redis.New(&cfg.Redis, logger)
	if err != nil {
		return nil, err
	}

	encoder := encoder.New()

	render := render.New(cfg.TemplatesPath, logger)

	serviceURLShortener := url_shortener.New(logger, rds, postgres)

	httpServer, err := ports.NewServer(&cfg.Http, logger, serviceURLShortener, encoder, render)
	if err != nil {
		return nil, err
	}

	return &Components{
		Postgres:   postgres,
		Redis:      rds,
		HttpServer: httpServer,
	}, nil
}

func (c *Components) Shutdown() {
	c.HttpServer.Stop()
	c.Postgres.CloseConnection()
	c.Redis.Close()
}

func SetupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = slogpretty.SetupPrettySlog()
	case envDev:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return logger
}
