package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"strconv"
	"url-shortner/internal/config"
	"url-shortner/internal/domain"
)

const Key = "links"

var errKeyDoesNotExists = errors.New("key does not exists")

type Redis struct {
	client redis.UniversalClient
	logger *slog.Logger
}

func New(config *config.RedisConfig, logger *slog.Logger) (*Redis, error) {
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    config.Hosts,
		Password: config.Password,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("storage.redis.New: %w", err)
	}

	return &Redis{
		client: client,
		logger: logger,
	}, nil
}

func (r *Redis) Close() {
	err := r.client.Close()
	if err != nil {
		r.logger.Error("storage.redis.Close", slog.String("error", err.Error()))
	}
}

func (r *Redis) QueryLinkByID(ctx context.Context, id int) (string, error) {
	if !r.hashExists(ctx) {
		return "", errKeyDoesNotExists
	}

	return r.client.HGet(ctx, Key, strconv.Itoa(id)).Result()
}

func (r *Redis) StoreLink(ctx context.Context, link *domain.Link) error {
	err := r.client.HSet(ctx, Key, strconv.Itoa(link.ID), link.URL).Err()
	if err != nil {
		return fmt.Errorf("storage.redis.AddRoomClient: %w", err)
	}

	return nil
}

func (r *Redis) hashExists(ctx context.Context) bool {
	exists, err := r.client.Exists(ctx, Key).Result()

	if err != nil || exists == 0 {
		return false
	}

	return true
}
