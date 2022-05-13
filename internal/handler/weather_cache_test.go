package handler

import (
	"context"
	"testing"

	"github.com/taadis/weather-api/internal/cache"
)

var (
	weatherCache *WeatherCache
)

func TestMain(m *testing.M) {
	ctx = context.Background()
	weatherCache = NewWeatherCache()
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
