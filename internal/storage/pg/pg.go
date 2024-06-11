package pg

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"url-shortner/internal/domain"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func New(pgURL string) (*Postgres, error) {
	config, err := pgxpool.ParseConfig(pgURL)
	if err != nil {
		return nil, fmt.Errorf("storage.pg.New: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("storage.pg.New: %w", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("storage.pg.New: %w", err)
	}

	return &Postgres{
		pool: pool,
	}, nil
}

func (pg *Postgres) CloseConnection() {
	pg.pool.Close()
}

func (pg *Postgres) GetByID(ctx context.Context, id int) (*domain.Link, error) {
	var link domain.Link
	err := pg.pool.QueryRow(ctx, "SELECT id, url FROM links WHERE id = $1", id).Scan(&link.ID, &link.URL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrURLNotFound
		}

		return nil, err
	}

	return &link, nil
}

func (pg *Postgres) GetByURL(ctx context.Context, url string) (*domain.Link, error) {
	var link domain.Link
	err := pg.pool.QueryRow(ctx, "SELECT id, url FROM links WHERE url = $1", url).Scan(&link.ID, &link.URL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrURLNotFound
		}

		return nil, err
	}

	return &link, nil
}

func (pg *Postgres) PersistURL(ctx context.Context, url string) (*domain.Link, error) {
	var newID int
	err := pg.pool.QueryRow(ctx, "INSERT INTO links (url) VALUES($1) returning id", url).Scan(&newID)
	if err != nil {
		return nil, err
	}

	return &domain.Link{
		ID:  newID,
		URL: url,
	}, nil
}
