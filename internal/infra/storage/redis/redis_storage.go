package redis

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/vaberof/hezzl-backend/internal/infra/storage"
	"time"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(client *redis.Client) *RedisStorage {
	return &RedisStorage{client: client}
}

func (rs *RedisStorage) Set(key, value string, exp time.Duration) error {
	err := rs.client.Set(context.Background(), key, value, exp).Err()
	if err != nil {
		return err
	}
	return nil
}

func (rs *RedisStorage) Get(key string) (string, error) {
	val, err := rs.client.Get(context.Background(), key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", storage.ErrRedisKeyNotFound
		}
		return "", err
	}
	return val, nil
}

func (rs *RedisStorage) Delete(key ...string) error {
	_, err := rs.client.Del(context.Background(), key...).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return storage.ErrRedisKeyNotFound
		}
		return err
	}
	return nil
}
