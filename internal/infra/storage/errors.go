package storage

import "errors"

var (
	ErrPostgresGoodNotFound = errors.New("good not found")

	ErrRedisKeyNotFound = errors.New("key not found")
)
