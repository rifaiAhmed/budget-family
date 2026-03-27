package config

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedis(cfg Config) (*redis.Client, error) {
	if !cfg.Redis.Enabled {
		return redis.NewClient(&redis.Options{Addr: cfg.Redis.Addr}), nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Pass,
		DB:       cfg.Redis.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
