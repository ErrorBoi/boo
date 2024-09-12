package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/errorboi/boo/internal/config"
	"github.com/errorboi/boo/store"
	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

func New(config config.RedisConfig) *RedisStore {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Address,
		Password: config.Password,
		DB:       config.DB,
	})

	return &RedisStore{
		client: rdb,
	}
}

func (s *RedisStore) CreateWithTTL(ctx context.Context, key string, val interface{}, ttl time.Duration) error {
	key = getNotifyKey(key)
	return s.client.Set(ctx, key, val, ttl).Err()
}

func (s *RedisStore) Get(ctx context.Context, key string) (string, error) {
	key = getNotifyKey(key)
	val, err := s.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", store.ErrEntityNotFound
		}
		return "", err
	}

	return val, nil
}

func (s *RedisStore) Delete(ctx context.Context, key string) error {
	key = getNotifyKey(key)
	err := s.client.Del(ctx, key).Err()

	return err
}

func getNotifyKey(key string) string {
	return fmt.Sprintf("notify:%s", key)
}
