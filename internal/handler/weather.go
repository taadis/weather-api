package handler

import (
	"context"
	"encoding/json"

	qweather "github.com/Ink-33/go-heweather/v7"
	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/util/log"
	qweathersdk "github.com/taadis/qweather-sdk-go"
	"github.com/taadis/weather-api/internal/conf"
	"github.com/taadis/weather-api/internal/model"
)

const (
	// 免费开发版为false，商业共享版与商业高性能版均为true
	isBusiness = false
)

var errJsonUnmarshal = errors.InternalServerError("", "序列化失败,请重试")

type Weather struct {
	qweatherCredential *qweather.Credential
	qweatherClient     *qweathersdk.Client
}

func NewWeather() *Weather {
	h := new(Weather)
	h.qweatherCredential = qweather.NewCredential(conf.GetPublicId(), conf.GetKey(), isBusiness) // 创建一个安全凭证
	h.qweatherClient = qweathersdk.NewClient()
	return h
}

func (h *Weather) TopCity(ctx context.Context, req *model.TopCityRequest, resp *model.TopCityResponse) error {
	v2TopCityReq := qweathersdk.NewV2TopCityRequest()
	v2TopCityReq.Key = conf.GetKey()
	v2TopCityResp, err := h.qweatherClient.V2TopCity(v2TopCityReq)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(v2TopCityResp.String()), resp)
	if err != nil {
		log.Errorf("TopCity json.Unmarshal error:%v, result:%s", err, v2TopCityResp.String())
		return errJsonUnmarshal
	}

	return nil
}

func (h *Weather) LookupCity(ctx context.Context, req *model.LookupCityRequest, resp *model.LookupCityResponse) error {
	v2LookupCityReq := qweathersdk.NewV2LookupCityRequest()
	v2LookupCityReq.Key = conf.GetKey()
	v2LookupCityReq.Location = req.Location
	v2LookupCityReq.Adm = req.Adm
	v2LookupCityResp, err := h.qweatherClient.V2LookupCity(v2LookupCityReq)
	if err != nil {
		log.Errorf("got V2LookupCity error: %v", err)
		return err
	}

	err = json.Unmarshal([]byte(v2LookupCityResp.String()), resp)
	if err != nil {
		return errJsonUnmarshal
	}

	return nil
}

// Indices 天气生活指数
func (h *Weather) Indices(ctx context.Context, req *model.WeatherIndicesRequest, resp *model.WeatherIndicesResponse) error {
	v7IndicesReq := qweathersdk.NewV7IndicesRequest()
	v7IndicesReq.Key = conf.GetKey()
	v7IndicesReq.IsDev = true
	v7IndicesReq.Location = req.Location
	v7IndicesReq.Type = req.Type
	v7IndicesReq.Duration = req.Duration
	v7IndicesResp, err := h.qweatherClient.V7Indices(v7IndicesReq)
	if err != nil {
		log.Errorf("got V7Indices error:%v, req:%v", err, v7IndicesReq)
		return err
	}

	err = json.Unmarshal([]byte(v7IndicesResp.String()), resp)
	if err != nil {
		return errJsonUnmarshal
	}

	return nil
}

// Now 实时天气
func (h *Weather) Now(ctx context.Context, req *model.WeatherNowRequest, resp *model.WeatherNowResponse) error {
	v7WeatherNowReq := qweathersdk.NewV7WeatherNowRequest()
	v7WeatherNowReq.Key = conf.GetKey()
	v7WeatherNowReq.IsDev = true
	v7WeatherNowReq.Location = req.Location
	v7WeatherNowResp, err := h.qweatherClient.V7WeatherNow(v7WeatherNowReq)
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
func (h *Weather) Forecast(ctx context.Context, req *model.WeatherForecastRequest, resp *model.WeatherForecastResponse) error {
	v7WeatherDaysReq := qweathersdk.NewV7WeatherDaysRequest()
	v7WeatherDaysReq.Key = conf.GetKey()
	v7WeatherDaysReq.IsDev = true
	v7WeatherDaysReq.Location = req.Location
	v7WeatherDaysReq.Duration = req.Duration
	v7WeatherDaysResp, err := h.qweatherClient.V7WeatherDays(v7WeatherDaysReq)
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
