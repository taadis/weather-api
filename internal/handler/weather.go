package handler

import (
	"context"
	"encoding/json"

	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/util/log"
	weatherSdk "github.com/taadis/qweather-sdk-go"
	"github.com/taadis/weather-api/internal/conf"
	"github.com/taadis/weather-api/internal/model"
)

var errJsonUnmarshal = errors.InternalServerError("", "序列化失败,请重试")

type Weather struct {
	weatherClient *weatherSdk.Client
	weatherCache  IWeatherCache
}

func NewWeather() *Weather {
	h := new(Weather)
	h.weatherClient = weatherSdk.NewClient()
	h.weatherCache = NewWeatherCache()
	return h
}

func (h *Weather) TopCity(ctx context.Context, _ *model.TopCityRequest, resp *model.TopCityResponse) error {
	s, err := h.weatherCache.GetSetTopCity(ctx)
	if err != nil {
		log.Errorf("TopCity cache.GetSetTopCity error:%+v", err)
		return err
	}

	err = json.Unmarshal([]byte(s), resp)
	if err != nil {
		log.Errorf("TopCity json.Unmarshal error:%+v, s:%s", err, s)
		return errJsonUnmarshal
	}

	return nil
}

func (h *Weather) LookupCity(ctx context.Context, req *model.LookupCityRequest, resp *model.LookupCityResponse) error {
	key := new(LookupCityKey)
	key.Location = req.Location
	key.Adm = req.Adm
	s, err := h.weatherCache.GetSetLookupCity(ctx, key)
	if err != nil {
		log.Errorf("LookupCity cache.GetSetLookupCity error:%+v", err)
		return err
	}

	err = json.Unmarshal([]byte(s), resp)
	if err != nil {
		log.Errorf("LookupCity json.Unmarshal error:%+v", err)
		return errJsonUnmarshal
	}

	return nil
}

// Indices 天气生活指数
func (h *Weather) Indices(ctx context.Context, req *model.WeatherIndicesRequest, resp *model.WeatherIndicesResponse) error {
	key := new(IndicesKey)
	key.Location = req.Location
	key.Type = req.Type
	key.Duration = req.Duration
	s, err := h.weatherCache.GetSetIndices(ctx, key)
	if err != nil {
		log.Errorf("Indices cache.GetSetIndices error:%+v", err)
		return err
	}

	err = json.Unmarshal([]byte(s), resp)
	if err != nil {
		log.Errorf("Indices json.Unmarshal error:%+v", err)
		return errJsonUnmarshal
	}

	return nil
}

// Now 实时天气
func (h *Weather) Now(_ context.Context, req *model.WeatherNowRequest, resp *model.WeatherNowResponse) error {
	v7WeatherNowReq := weatherSdk.NewV7WeatherNowRequest()
	v7WeatherNowReq.Key = conf.GetKey()
	v7WeatherNowReq.IsDev = true
	v7WeatherNowReq.Location = req.Location
	v7WeatherNowResp, err := h.weatherClient.V7WeatherNow(v7WeatherNowReq)
	if err != nil {
		log.Errorf("got V7WeatherNow error:%v, req:%v", err, v7WeatherNowReq)
		return err
	}

	err = json.Unmarshal([]byte(v7WeatherNowResp.String()), resp)
	if err != nil {
		return errJsonUnmarshal
	}

	return nil
}

// Forecast 天气预报
func (h *Weather) Forecast(_ context.Context, req *model.WeatherForecastRequest, resp *model.WeatherForecastResponse) error {
	v7WeatherDaysReq := weatherSdk.NewV7WeatherDaysRequest()
	v7WeatherDaysReq.Key = conf.GetKey()
	v7WeatherDaysReq.IsDev = true
	v7WeatherDaysReq.Location = req.Location
	v7WeatherDaysReq.Duration = req.Duration
	v7WeatherDaysResp, err := h.weatherClient.V7WeatherDays(v7WeatherDaysReq)
	if err != nil {
		log.Errorf("got V7WeatherDays error:%v, req:%v", err, v7WeatherDaysReq)
		return err
	}

	err = json.Unmarshal([]byte(v7WeatherDaysResp.String()), resp)
	if err != nil {
		return errJsonUnmarshal
	}

	return nil
}
