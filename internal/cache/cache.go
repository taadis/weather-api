package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cache struct {
	rdb *redis.Client
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

func NewCache(rdb *redis.Client) *Cache {
	c := new(Cache)
	c.rdb = rdb
	return c
}

func (c *Cache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.rdb.Set(ctx, key, value, expiration).Err()
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	return c.rdb.Get(ctx, key).Result()
}

func (c *Cache) Del(ctx context.Context, key string) error {
	return c.rdb.Del(ctx, key).Err()
}
