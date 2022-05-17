package handler

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/taadis/qweather-sdk-go"
	"github.com/taadis/weather-api/internal/cache"
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
	weatherCache = NewWeatherCache(cache, weatherClient)
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
