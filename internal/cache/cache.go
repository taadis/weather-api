package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/micro/go-micro/util/log"
	qWeather "github.com/taadis/qweather-sdk-go"
	"github.com/taadis/weather-api/internal/conf"
)

const (
	keyTopCity = "top-city"
)

type Cache struct {
	rdb *redis.Client
}

func IsNil(err error) bool {
	if err == redis.Nil {
		return true
	}

	return false
}

func NewCache(rdb *redis.Client) *Cache {
	c := new(Cache)
	c.rdb = rdb
	return c
}

func (c *Cache) SetTopCity(ctx context.Context, value string) error {
	err := c.rdb.Set(ctx, keyTopCity, value, 10*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) GetTopCity(ctx context.Context) (string, error) {
	s, err := c.rdb.Get(ctx, keyTopCity).Result()
	if err != nil {
		return "", err
	}

	return s, nil
}

func (c *Cache) GetTopCityWithSet(ctx context.Context) (string, error) {
	s, err := c.GetTopCity(ctx)
	if err != nil {
		if IsNil(err) {
			log.Logf("cache.GetTopCity is nil")
			v2TopCityReq := qWeather.NewV2TopCityRequest()
			v2TopCityReq.Key = conf.GetKey()
			v2TopCityResp, err := qWeather.NewClient().V2TopCity(v2TopCityReq)
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
	}
	return s, nil
}
