package ports

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"log/slog"
	"net/http"
	"time"
	"url-shortner/internal/config"
	"url-shortner/internal/ports/rest"
	mwlogger "url-shortner/pkg/logger/middleware"
	"url-shortner/pkg/rate_limiter"
)

type Server struct {
	logger          *slog.Logger
	server          *http.Server
	shutDownTimeout time.Duration
}

func NewServer(config *config.HTTPConfig, logger *slog.Logger, serviceURLShortener rest.ServiceURLShortener, serviceEncoder rest.ServiceEncoder, serviceRender rest.ServiceRender) (*Server, error) {
	httpHandler := rest.NewHandler(logger, serviceURLShortener, serviceEncoder, serviceRender)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", config.Port),
		Handler:      InitRouter(httpHandler, logger, &config.Limiter),
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	}

	return &Server{
		server:          server,
		shutDownTimeout: config.ShutdownTimeout,
		logger:          logger,
	}, nil
}

func InitRouter(handler *rest.Handler, logger *slog.Logger, limiter *config.Limiter) *chi.Mux {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300, // максимальный срок кэширования предварительных запросов
	}))

	mux.Use(rate_limiter.Limit(limiter.RPS, limiter.Burst, limiter.TTL, logger))
	mux.Use(middleware.Recoverer)
	mux.Use(mwlogger.Log(logger))

	mux.Get("/", handler.Homepage)
	mux.Get("/favicon.ico", handler.Icon)
	mux.Get("/{code}", handler.ProxyURLCode)
	mux.Post("/api/urls", handler.RegisterURL)

	return mux
}

func (s *Server) Run(ctx context.Context) error {
	errResult := make(chan error)
	go func() {
		s.logger.Info(fmt.Sprintf("starting listening: %s", s.server.Addr))

		errResult <- s.server.ListenAndServe()
	}()

	var err error
	select {
	case <-ctx.Done():
		return ctx.Err()

	case err = <-errResult:
	}
	return err
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutDownTimeout)
	defer cancel()
	err := s.server.Shutdown(ctx)
	if err != nil {
		s.logger.Error("failed to shutdown HTTP Server", slog.String("error", err.Error()))
	}
}
