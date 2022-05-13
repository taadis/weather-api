package cache

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
)

var (
	ctx         context.Context
	cacheClient *Cache
)

func TestMain(m *testing.M) {
	ctx = context.Background()
	mockRdb, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     mockRdb.Addr(),
		Password: "",
		DB:       0,
	})
	cacheClient = NewCache(rdb)
	m.Run()
	rdb.Close()
}

func TestCache_Set(t *testing.T) {
	err := cacheClient.Set(ctx, "key1", "value1", 3*time.Second)
	if err != nil {
		t.Logf("cache set error:%+v", err)
	}
}

func TestCache_Get(t *testing.T) {
	s, err := cacheClient.Get(ctx, "key1")
	if err != nil {
		if IsNil(err) {
			t.Logf("cache get key1 is nil")
			return
		}
		t.Fatalf("cache get key1 error:%+v", err)
	}

	t.Logf("cache get key1:%s", s)
}

func TestCache_Delete(t *testing.T) {
	err := cacheClient.Del(ctx, "key1")
	if err != nil {
		t.Fatalf("cache delete key1 error:%+v", err)
	}
}

func TestCache(t *testing.T) {
	key := "key1"
	t.Run("KeyNotFound", func(t *testing.T) {
		_, err := cacheClient.Get(ctx, key)
		if !IsNil(err) {
			t.Fatal(err)
		}
	})
	t.Run("Set", func(t *testing.T) {
		err := cacheClient.Set(ctx, key, "value1", time.Second)
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("Get", func(t *testing.T) {
		value, err := cacheClient.Get(ctx, key)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("cache get %s=%s", key, value)
	})
	t.Run("Delete", func(t *testing.T) {
		err := cacheClient.Del(ctx, key)
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("Deleted", func(t *testing.T) {
		_, err := cacheClient.Get(ctx, key)
		if !IsNil(err) {
			t.Fatal(err)
		}
	})
}
