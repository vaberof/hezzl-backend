package good

import (
	"time"
)

type InMemoryStorage interface {
	Set(key, value string, exp time.Duration) error
	Get(key string) (string, error)
}
