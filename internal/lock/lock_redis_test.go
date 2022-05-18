package lock

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
)

func TestLockRedis(t *testing.T) {
	mockRedis := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{
		Addr:     mockRedis.Addr(),
		Password: "",
		DB:       0,
	})
	locker := NewLockRedis(rdb)
	key := "test:lock"
	ctx := context.Background()
	t.Run("Lock", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			result, err := locker.Lock(ctx, key, 2*time.Second)
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("got lock result:%v", result)
		}
		err := locker.Unlock(ctx, key)
		if err != nil {
			t.Fatal(err)
		}
		result, err := locker.Lock(ctx, key, 2*time.Second)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("got lock result:%v", result)
	})
}
