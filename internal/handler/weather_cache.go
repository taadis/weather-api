package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/micro/go-micro/util/log"
	weatherSdk "github.com/taadis/qweather-sdk-go"
	"github.com/taadis/weather-api/internal/cache"
	"github.com/taadis/weather-api/internal/conf"
	"github.com/taadis/weather-api/internal/lock"
)

const (
	defaultExpiration = 10 * time.Minute
	KeyTopCity        = "top-city"
	KeyLookupCity     = "lookup-city-location-%s-adm-%s"
	KeyIndices        = "indices-location-%s-type-%s-duration-%s"
	KeyNow            = "now-location-%s"
	KeyForecast       = "forecast-location-%s-duration-%s"
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
	SetNow(ctx context.Context, key string, value string) error
	GetNow(ctx context.Context, key string) (string, error)
	GetSetNow(ctx context.Context, key *NowKey) (string, error)
	SetForecast(ctx context.Context, key string, value string) error
	GetForecast(ctx context.Context, key string) (string, error)
	GetSetForecast(ctx context.Context, key *ForecastKey) (string, error)
}

type WeatherCache struct {
	cache         cache.ICache
	weatherClient *weatherSdk.Client
	locker        lock.Locker
}

func NewWeatherCache(cache cache.ICache, weatherClient *weatherSdk.Client, locker lock.Locker) *WeatherCache {
	c := new(WeatherCache)
	c.cache = cache
	c.weatherClient = weatherClient
	c.locker = locker
	return c
}

func (c *WeatherCache) SetTopCity(ctx context.Context, value string) error {
	return c.cache.Set(ctx, KeyTopCity, value, defaultExpiration)
}

func (c *WeatherCache) GetTopCity(ctx context.Context) (string, error) {
	return c.cache.Get(ctx, KeyTopCity)
}

func (c *WeatherCache) GetSetTopCity(ctx context.Context) (string, error) {
	s, err := c.GetTopCity(ctx)
	if err != nil {
		if cache.IsNil(err) {
			// wait lock
			keyLockTopCity := "lock:top:city"
			err := c.locker.LockWait(ctx, keyLockTopCity, 2*time.Second)
			if err != nil {
				log.Errorf("locker.LockWait error:%+v, key=%s", err, keyLockTopCity)
				return "", err
			}
			defer func() {
				_ = c.locker.Unlock(ctx, keyLockTopCity)
			}()

			// get again
			s, err := c.GetTopCity(ctx)
			if err == nil && !cache.IsNil(err) {
				return s, nil
			}

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
			// wait lock
			keyLockLookupCity := fmt.Sprintf("lock:%s", ky)
			err := c.locker.LockWait(ctx, keyLockLookupCity, 2*time.Second)
			if err != nil {
				log.Errorf("locker.LockWait error:%+v, key=%s", err, keyLockLookupCity)
				return "", err
			}
			defer func() {
				_ = c.locker.Unlock(ctx, keyLockLookupCity)
			}()

			// 锁内再读一次
			s, err := c.GetLookupCity(ctx, ky)
			if err == nil && !cache.IsNil(err) {
				return s, nil
			}

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
	s, err := c.GetIndices(ctx, ky)
	if err != nil {
		if cache.IsNil(err) {
			// wait lock
			keyLockIndices := fmt.Sprintf("lock:%s", ky)
			err := c.locker.LockWait(ctx, keyLockIndices, 2*time.Second)
			if err != nil {
				log.Errorf("locker.LockWait error:%+v, key=%s", keyLockIndices)
				return "", err
			}
			defer func() {
				_ = c.locker.Unlock(ctx, keyLockIndices)
			}()

			// get again
			s, err := c.GetIndices(ctx, ky)
			if !cache.IsNil(err) && err == nil {
				return s, nil
			}

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

type NowKey struct {
	Location string
}

func (c *WeatherCache) SetNow(ctx context.Context, key string, value string) error {
	return c.cache.Set(ctx, key, value, defaultExpiration)
}

func (c *WeatherCache) GetNow(ctx context.Context, key string) (string, error) {
	return c.cache.Get(ctx, key)
}

func (c *WeatherCache) GetSetNow(ctx context.Context, key *NowKey) (string, error) {
	ky := cache.BuildKey(KeyNow, key.Location)
	s, err := c.GetNow(ctx, ky)
	if err != nil {
		if cache.IsNil(err) {
			// wait lock
			keyLockNow := fmt.Sprintf("lock:%s", ky)
			err := c.locker.LockWait(ctx, keyLockNow, 2*time.Second)
			if err != nil {
				log.Errorf("locker.LockWait error:%+v, key=%s", err, keyLockNow)
				return "", err
			}
			defer func() {
				_ = c.locker.Unlock(ctx, keyLockNow)
			}()

			// get again
			s, err := c.GetNow(ctx, ky)
			if !cache.IsNil(err) && err == nil {
				return s, nil
			}

			log.Infof("cache.GetNow is nil,key=%s", ky)
			v7WeatherNowReq := weatherSdk.NewV7WeatherNowRequest()
			v7WeatherNowReq.Key = conf.GetKey()
			v7WeatherNowReq.IsDev = true
			v7WeatherNowReq.Location = key.Location
			v7WeatherNowResp, err := c.weatherClient.V7WeatherNow(v7WeatherNowReq)
			if err != nil {
				log.Errorf("got V7WeatherNow error:%v, req:%v", err, v7WeatherNowReq)
				return "", err
			}

			s = v7WeatherNowResp.String()
			err = c.SetNow(ctx, ky, s)
			if err != nil {
				log.Errorf("cache.SetNow error:%+v,key=%s,value=%s", err, ky, s)
				return "", err
			}

			return s, nil
		}

		log.Errorf("cache.GetNow error:%+v,key=%s", err, ky)
		return "", err
	}
	return s, nil
}

type ForecastKey struct {
	Location string
	Duration string
}

func (c *WeatherCache) SetForecast(ctx context.Context, key string, value string) error {
	return c.cache.Set(ctx, key, value, defaultExpiration)
}

func (c *WeatherCache) GetForecast(ctx context.Context, key string) (string, error) {
	return c.cache.Get(ctx, key)
}

func (c *WeatherCache) GetSetForecast(ctx context.Context, key *ForecastKey) (string, error) {
	ky := cache.BuildKey(KeyForecast, key.Location, key.Duration)
	s, err := c.GetForecast(ctx, ky)
	if err != nil {
		if cache.IsNil(err) {
			// wait lock
			keyLockForecast := fmt.Sprintf("lock:%s", ky)
			err := c.locker.LockWait(ctx, keyLockForecast, 2*time.Second)
			if err != nil {
				log.Errorf("locker.LockWait error:%+v, key=%s", err, keyLockForecast)
				return "", err
			}
			defer func() {
				_ = c.locker.Unlock(ctx, keyLockForecast)
			}()

			// get again
			s, err := c.GetForecast(ctx, ky)
			if !cache.IsNil(err) && err == nil {
				return s, nil
			}

			log.Infof("cache.GetForecast is nil,key=%s", ky)
			v7WeatherDaysReq := weatherSdk.NewV7WeatherDaysRequest()
			v7WeatherDaysReq.Key = conf.GetKey()
			v7WeatherDaysReq.IsDev = true
			v7WeatherDaysReq.Location = key.Location
			v7WeatherDaysReq.Duration = key.Duration
			v7WeatherDaysResp, err := c.weatherClient.V7WeatherDays(v7WeatherDaysReq)
			if err != nil {
				log.Errorf("got V7WeatherDays error:%v, req:%v", err, v7WeatherDaysReq)
				return "", err
			}

			s = v7WeatherDaysResp.String()
			err = c.SetForecast(ctx, ky, s)
			if err != nil {
				log.Errorf("cache.SetForecast error:%+v,key=%s,value=%s", err, ky, s)
				return "", err
			}

			return s, nil
		}

		log.Errorf("cache.GetForecast error:%+v,key=%s", err, ky)
		return "", err
	}

	return s, nil
}
