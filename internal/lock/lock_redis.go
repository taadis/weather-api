package lock

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type lockRedis struct {
	rdb *redis.Client
}

func (l *lockRedis) Lock(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return l.rdb.SetNX(ctx, key, byte(1), expiration).Result()
}

func (l *lockRedis) LockWait(ctx context.Context, key string, expiration time.Duration) error {
	for {
		exists, err := l.rdb.Exists(ctx, key).Result()
		if err != nil {
			return err
		}

		// 不存在则设置加锁
		if exists == 0 {
			locked, err := l.Lock(ctx, key, expiration)
			if err != nil {
				return err
			}
			if locked {
				return nil
			}
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func (l *lockRedis) Unlock(ctx context.Context, key string) error {
	return l.rdb.Del(ctx, key).Err()
}

func NewLockRedis(rdb *redis.Client) Locker {
	l := new(lockRedis)
	l.rdb = rdb
	return l
}
