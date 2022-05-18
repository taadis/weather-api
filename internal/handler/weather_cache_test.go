package handler

import (
	"context"
	"sync"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/taadis/qweather-sdk-go"
	"github.com/taadis/weather-api/internal/cache"
	"github.com/taadis/weather-api/internal/lock"
)

var (
	weatherCache *WeatherCache
)

func TestMain(m *testing.M) {
	ctx = context.Background()
	mockRedis, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     mockRedis.Addr(),
		Password: "",
		DB:       0,
	})
	cache := cache.NewCache(rdb)
	weatherClient := qweather.NewClient()
	locker := lock.NewLockRedis(rdb)
	weatherCache = NewWeatherCache(cache, weatherClient, locker)
	weather = NewWeather(weatherCache)
	m.Run()
}

func TestCache_GetTopCity(t *testing.T) {
	s, err := weatherCache.GetTopCity(ctx)
	if err != nil {
		if cache.IsNil(err) {
			t.Logf("cacheClient.GetTopCity result is nil")
			return
		}
		t.Fatalf("cacheClient.GetTopCity error:%+v", err)
	}

	t.Logf("cacheClient.GetTopCity result:%s", s)
}

func TestCache_SetTopCity(t *testing.T) {
	value := "xxx"
	err := weatherCache.SetTopCity(ctx, value)
	if err != nil {
		t.Fatalf("cacheClient.SetTopCity error:%+v", err)
	}
}

func TestCache_GetTopCityWithSet(t *testing.T) {
	s, err := weatherCache.GetSetTopCity(ctx)
	if err != nil {
		t.Fatalf("cacheClient.GetTopCityWithSet error:%+v", err)
	}

	t.Logf("cacheClient.GetTopCityWithSet result:%s", s)
}

// 并发请求的时候,会产生缓存击穿问题,补充锁机制及重读
func TestWeatherCache_GetSetTopCity_Breakdown(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			s, err := weatherCache.GetSetTopCity(ctx)
			if err != nil {
				t.Fatalf("cacheClient.GetTopCityWithSet error:%+v", err)
			}

			t.Logf("cacheClient.GetTopCityWithSet result:%s", s)
		}(&wg)
	}
	wg.Wait()
}
