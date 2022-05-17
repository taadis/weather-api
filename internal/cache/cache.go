package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type ICache interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
}

type Cache struct {
	*redis.Client
}

func NewCache(rdb *redis.Client) *Cache {
	c := new(Cache)
	c.Client = rdb
	return c
}

func BuildKey(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}

func IsNil(err error) bool {
	if err == redis.Nil {
		return true
	}

	return false
}

func (c *Cache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.Client.Set(ctx, key, value, expiration).Err()
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	return c.Client.Get(ctx, key).Result()
}

func (c *Cache) Del(ctx context.Context, key string) error {
	return c.Client.Del(ctx, key).Err()
}
