package cache

import (
	"context"
	"testing"

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
}

func TestCache_GetTopCity(t *testing.T) {
	s, err := cacheClient.GetTopCity(ctx)
	if err != nil {
		if err == redis.Nil {
			t.Logf("cacheClient.GetTopCity result is nil")
			return
		}
		t.Fatalf("cacheClient.GetTopCity error:%+v", err)
	}

	t.Logf("cacheClient.GetTopCity result:%s", s)
}

func TestCache_SetTopCity(t *testing.T) {
	value := "xxx"
	err := cacheClient.SetTopCity(ctx, value)
	if err != nil {
		t.Fatalf("cacheClient.SetTopCity error:%+v", err)
	}
}

func TestCache_GetTopCityWithSet(t *testing.T) {
	s, err := cacheClient.GetTopCityWithSet(ctx)
	if err != nil {
		t.Fatalf("cacheClient.GetTopCityWithSet error:%+v", err)
	}

	t.Logf("cacheClient.GetTopCityWithSet result:%s", s)
}
