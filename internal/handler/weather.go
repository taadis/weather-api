package handler

import (
	"context"
	"encoding/json"

	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/util/log"
	"github.com/taadis/weather-api/internal/model"
)

var errJsonUnmarshal = errors.InternalServerError("", "序列化失败,请重试")

type IWeather interface {
	TopCity(ctx context.Context, _ *model.TopCityRequest, resp *model.TopCityResponse) error
	LookupCity(ctx context.Context, req *model.LookupCityRequest, resp *model.LookupCityResponse) error
	Indices(ctx context.Context, req *model.WeatherIndicesRequest, resp *model.WeatherIndicesResponse) error
	Now(ctx context.Context, req *model.WeatherNowRequest, resp *model.WeatherNowResponse) error
	Forecast(ctx context.Context, req *model.WeatherForecastRequest, resp *model.WeatherForecastResponse) error
}

type Weather struct {
	weatherCache IWeatherCache
}

func NewWeather(weatherCache IWeatherCache) IWeather {
	h := new(Weather)
	h.weatherCache = weatherCache
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
func (h *Weather) Now(ctx context.Context, req *model.WeatherNowRequest, resp *model.WeatherNowResponse) error {
	key := new(NowKey)
	key.Location = req.Location
	s, err := h.weatherCache.GetSetNow(ctx, key)
	if err != nil {
		log.Errorf("Now cache.GetSetNow error:%+v", err)
		return err
	}

	err = json.Unmarshal([]byte(s), resp)
	if err != nil {
		log.Errorf("Now json.Unmarshal error:%+v,s:%s", s)
		return errJsonUnmarshal
	}

	return nil
}

// Forecast 天气预报
func (h *Weather) Forecast(ctx context.Context, req *model.WeatherForecastRequest, resp *model.WeatherForecastResponse) error {
	key := new(ForecastKey)
	key.Location = req.Location
	key.Duration = req.Duration
	s, err := h.weatherCache.GetSetForecast(ctx, key)
	if err != nil {
		log.Errorf("Forecast cache.GetSetForecast error:%+v", err)
		return err
	}

	err = json.Unmarshal([]byte(s), resp)
	if err != nil {
		log.Errorf("Forecast json.Unmarshal error:%+v,s:%s", err, s)
		return errJsonUnmarshal
	}

	return nil
}
