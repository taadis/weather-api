package lock

import (
	"context"
	"time"
)

type Locker interface {
	Lock(ctx context.Context, key string, expiration time.Duration) (bool, error)
	LockWait(ctx context.Context, key string, expiration time.Duration) error
	Unlock(ctx context.Context, key string) error
}
