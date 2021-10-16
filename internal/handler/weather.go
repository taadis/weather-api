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
	client := qweather.NewGeoCityClient(req.Location)
	lookupCityResp, err := client.Run(h.qweatherCredential, nil)
	if err != nil {
		log.Errorf("qweather.NewGeoCityClient client.Run lookupCity error: %v", err)
		return err
	}

	err = json.Unmarshal([]byte(lookupCityResp), resp)
	if err != nil {
		return errJsonUnmarshal
	}

	return nil
}

// Indices 天气生活指数
func (h *Weather) Indices(ctx context.Context, req *model.WeatherIndicesRequest, resp *model.WeatherIndicesResponse) error {
	client, err := qweather.NewLiveIndexClient(req.Location, req.Type, req.Duration)
	if err != nil {
		log.Errorf("qweather.NewLiveIndexClient errror:%v, req:%v", err, req)
		return err
	}

	result, err := client.Run(h.qweatherCredential, &qweather.ClientConfig{
		Language: "cn",
	})
	if err != nil {
		log.Errorf("qweather.NewLiveIndexClient Run error:%v", err)
		return err
	}

	err = json.Unmarshal([]byte(result), resp)
	if err != nil {
		log.Errorf("Weather.Indices json.Unmarshal error:%v, result:%s", err, result)
		return errJsonUnmarshal
	}

	return nil
}

// Now 实时天气
func (h *Weather) Now(ctx context.Context, req *model.WeatherNowRequest, resp *model.WeatherNowResponse) error {
	client := qweather.NewRealTimeWeatherClient(req.Location)
	result, err := client.Run(h.qweatherCredential, nil)
	if err != nil {
		log.Errorf("qweather.NewRealTimeWeatherClient Run error: %v", err)
		return err
	}

	err = json.Unmarshal([]byte(result), resp)
	if err != nil {
		log.Errorf("Weather.Now json.Unmarshal error:%v, result:%s", err, result)
		return errJsonUnmarshal
	}

	return nil
}

// Forecast 天气预报
func (h *Weather) Forecast(ctx context.Context, req *model.WeatherForecastRequest, resp *model.WeatherForecastResponse) error {
	client, err := qweather.NewWeatherForecastClient(req.Location, req.Duration)
	if err != nil {
		log.Errorf("qweather.NewWeatherForecastClient error:%v, req:%+v", err, req)
		return err
	}

	result, err := client.Run(h.qweatherCredential, nil)
	if err != nil {
		log.Errorf("qweather.NewWeatherForecastClient Run error:%v, req:%+v", err, req)
		return err
	}

	err = json.Unmarshal([]byte(result), resp)
	if err != nil {
		log.Errorf("Weather.Forecast json.Unmarshal error:%v, result:%s", result)
		return errJsonUnmarshal
	}

	return nil
}
