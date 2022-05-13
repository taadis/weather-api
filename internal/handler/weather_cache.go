package handler

import (
	"context"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/micro/go-micro/util/log"
	weatherSdk "github.com/taadis/qweather-sdk-go"
	"github.com/taadis/weather-api/internal/cache"
	"github.com/taadis/weather-api/internal/conf"
)

const (
	defaultExpiration = 10 * time.Minute
	KeyTopCity        = "top-city"
	KeyLookupCity     = "lookup-city-location-%s-adm-%s"
	KeyIndices        = "indices-location-%s-type-%s-duration-%s"
)

type IWeatherCache interface {
	SetTopCity(ctx context.Context, value string) error
	GetTopCity(ctx context.Context) (string, error)
	GetSetTopCity(ctx context.Context) (string, error)
	SetLookupCity(ctx context.Context, key string, value string) error
	GetLookupCity(ctx context.Context, key string) (string, error)
	GetSetLookupCity(ctx context.Context, key *LookupCityKey) (string, error)
	SetIndices(ctx context.Context, key string, value string) error
	GetIndices(ctx context.Context, key string) (string, error)
	GetSetIndices(ctx context.Context, key *IndicesKey) (string, error)
}

type WeatherCache struct {
	cache         *cache.Cache
	weatherClient *weatherSdk.Client
}

func NewWeatherCache() *WeatherCache {
	c := new(WeatherCache)
	c.weatherClient = weatherSdk.NewClient()
	mockRedis, _ := miniredis.Run()
	rdbOptions := &redis.Options{
		Addr:     mockRedis.Addr(),
		Password: "",
		DB:       0,
	}
	c.cache = cache.NewCache(redis.NewClient(rdbOptions))
	return c
}

func (c *WeatherCache) SetTopCity(ctx context.Context, value string) error {
	return c.cache.Set(ctx, KeyLookupCity, value, 10*time.Minute)
}

func (c *WeatherCache) GetTopCity(ctx context.Context) (string, error) {
	return c.cache.Get(ctx, KeyTopCity)
}

func (c *WeatherCache) GetSetTopCity(ctx context.Context) (string, error) {
	s, err := c.GetTopCity(ctx)
	if err != nil {
		if cache.IsNil(err) {
			log.Logf("cache.GetTopCity is nil")
			v2TopCityReq := weatherSdk.NewV2TopCityRequest()
			v2TopCityReq.Key = conf.GetKey()
			v2TopCityResp, err := weatherSdk.NewClient().V2TopCity(v2TopCityReq)
			if err != nil {
				log.Errorf("V2TopCity error:%+v", err)
				return "", err
			}

			s = v2TopCityResp.String()
			err = c.SetTopCity(ctx, s)
			if err != nil {
				return "", err
			}
			return s, nil
		}
		log.Errorf("cache.GetTopCity error:%+v", err)
		return "", err
	}
	return s, nil
}

type LookupCityKey struct {
	Location string
	Adm      string
}

func (c *WeatherCache) SetLookupCity(ctx context.Context, key string, value string) error {
	return c.cache.Set(ctx, key, value, defaultExpiration)
}

func (c *WeatherCache) GetLookupCity(ctx context.Context, key string) (string, error) {
	return c.cache.Get(ctx, key)
}

func (c *WeatherCache) GetSetLookupCity(ctx context.Context, key *LookupCityKey) (string, error) {
	ky := cache.BuildKey(KeyLookupCity, key.Location, key.Adm)
	s, err := c.GetLookupCity(ctx, ky)
	if err != nil {
		if cache.IsNil(err) {
			log.Logf("cache.GetLookupCity is nil, key=%s", ky)
			v2LookupCityReq := weatherSdk.NewV2LookupCityRequest()
			v2LookupCityReq.Key = conf.GetKey()
			v2LookupCityReq.Location = key.Location
			v2LookupCityReq.Adm = key.Adm
			v2LookupCityResp, err := c.weatherClient.V2LookupCity(v2LookupCityReq)
			if err != nil {
				log.Errorf("got V2LookupCity error: %+v", err)
				return "", err
			}

			err = c.SetLookupCity(ctx, ky, v2LookupCityResp.String())
			if err != nil {
				log.Errorf("cache.SetLookupCity error:%+v,key=%s,value=%s", err, ky, v2LookupCityResp.String())
				return "", err
			}

			return v2LookupCityResp.String(), nil
		}
		log.Errorf("cache.GetLookupCity error:%+v,key=%s", err, ky)
		return "", err
	}
	return s, nil
}

func (c *WeatherCache) SetIndices(ctx context.Context, key string, value string) error {
	return c.cache.Set(ctx, key, value, defaultExpiration)
}

func (c *WeatherCache) GetIndices(ctx context.Context, key string) (string, error) {
	return c.cache.Get(ctx, key)
}

type IndicesKey struct {
	Location string
	Type     string
	Duration string
}

func (c *WeatherCache) GetSetIndices(ctx context.Context, key *IndicesKey) (string, error) {
	ky := cache.BuildKey(KeyIndices, key.Location, key.Type, key.Duration)
	s, err := c.cache.Get(ctx, ky)
	if err != nil {
		if cache.IsNil(err) {
			log.Logf("cache.GetIndices is nil,key=%s", ky)
			v7IndicesReq := weatherSdk.NewV7IndicesRequest()
			v7IndicesReq.Key = conf.GetKey()
			v7IndicesReq.IsDev = true
			v7IndicesReq.Location = key.Location
			v7IndicesReq.Type = key.Type
			v7IndicesReq.Duration = key.Duration
			v7IndicesResp, err := c.weatherClient.V7Indices(v7IndicesReq)
			if err != nil {
				log.Errorf("got V7Indices error:%v, req:%v", err, v7IndicesReq)
				return "", err
			}

			s = v7IndicesResp.String()

			err = c.SetIndices(ctx, ky, s)
			if err != nil {
				log.Errorf("cache.SetIndices error:%+v,key=%s,value=%s", err, ky, s)
				return "", err
			}

			return s, nil
		}
		log.Errorf("cache.GetIndices error:%+v,key=%s", ky)
		return "", err
	}
	return s, nil
}
