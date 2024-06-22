package url_shortener

import (
	"context"
	"errors"
	"log/slog"
	"url-shortner/internal/domain"
)

type Cache interface {
	QueryLinkByID(ctx context.Context, id int) (string, error)
	StoreLink(ctx context.Context, link *domain.Link) error
}

type DB interface {
	GetByID(ctx context.Context, id int) (*domain.Link, error)
	GetByURL(ctx context.Context, url string) (*domain.Link, error)
	PersistURL(ctx context.Context, url string) (*domain.Link, error)
}

type URLShortener struct {
	logger *slog.Logger
	cache  Cache
	db     DB
}

func New(logger *slog.Logger, cache Cache, db DB) *URLShortener {
	return &URLShortener{
		logger: logger,
		cache:  cache,
		db:     db,
	}
}

func (u *URLShortener) Proxy(ctx context.Context, id int) (string, error) {
	// first check if the link exists in Redis
	redisLink, err := u.cache.QueryLinkByID(ctx, id)
	if err == nil {
		return redisLink, nil
	}

	// link not found on Redis.
	// So, let's query the DB
	dbLink, err := u.db.GetByID(ctx, id)
	if err != nil {
		return "", err
	}

	// store the link on Redis
	err = u.cache.StoreLink(ctx, dbLink)
	if err != nil {
		u.logger.Error("cache error", slog.String("message", err.Error()))
	}

	return dbLink.URL, nil
}

func (u *URLShortener) Create(ctx context.Context, url string) (*domain.Link, error) {
	// check if link already exists on database
	storedLink, err := u.db.GetByURL(ctx, url)
	if err == nil {
		return storedLink, nil
	}

	if !errors.Is(err, domain.ErrURLNotFound) {
		return nil, err
	}

	// It's a new link, so let's persist it
	newLink, err := u.db.PersistURL(ctx, url)
	if err != nil {
		return nil, err
	}

	return newLink, nil
}
