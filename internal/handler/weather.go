package handler

import (
	"context"
	"log"
	"os"

	qweather "github.com/Ink-33/go-heweather/v7"
	"github.com/taadis/weather-api/proto"
)

const (
	QWEATHER_PUBLIC_ID = "QWEATHER_PUBLIC_ID"
	QWEATHER_KEY       = "QWEATHER_KEY"
	isBusiness         = false // 免费开发版为false，商业共享版与商业高性能版均为true
)

type WeatherHandler struct {
	qweatherCredential *qweather.Credential
}

func NewWeatherHandler() proto.WeatherHandler {
	h := new(WeatherHandler)
	h.qweatherCredential = qweather.NewCredential(h.getPublicId(), h.getKey(), isBusiness) // 创建一个安全凭证
	return h
}

func (h *WeatherHandler) getPublicId() string {
	publicId := os.Getenv(QWEATHER_PUBLIC_ID)
	return publicId
}

func (h *WeatherHandler) getKey() string {
	key := os.Getenv(QWEATHER_KEY)
	return key
}

// Indices 天气生活指数
func (h *WeatherHandler) Indices(ctx context.Context, req *proto.IndicesRequest, resp *proto.IndicesResponse) error {
	client, err := qweather.NewLiveIndexClient(req.Location, req.Type, req.Duration)
	if err != nil {
		log.Printf("qweather.NewLiveIndexClient errror:%v, req:%v", err, req)
		return err
	}

	result, err := client.Run(h.qweatherCredential, &qweather.ClientConfig{
		Language: "cn",
	})
	if err != nil {
		log.Printf("qweather.NewLiveIndexClient Run error:%v", err)
		return err
	}

	resp.Result = result
	return nil
}

// Now 实时天气
func (h *WeatherHandler) Now(ctx context.Context, req *proto.NowRequest, resp *proto.NowResponse) error {
	client := qweather.NewRealTimeWeatherClient(req.Location)
	result, err := client.Run(h.qweatherCredential, nil)
	if err != nil {
		log.Printf("qweather.NewRealTimeWeatherClient Run error: %v", err)
		return err
	}

	resp.Result = result
	return nil
}

// Forecast 天气预报
func (h *WeatherHandler) Forecast(ctx context.Context, req *proto.ForecastRequest, resp *proto.ForecastResponse) error {
	client, err := qweather.NewWeatherForecastClient(req.Location, req.Duration)
	if err != nil {
		log.Printf("qweather.NewWeatherForecastClient error:%v, req:%+v", err, req)
		return err
	}

	result, err := client.Run(h.qweatherCredential, nil)
	if err != nil {
		log.Printf("qweather.NewWeatherForecastClient Run error:%v, req:%+v", err, req)
		return err
	}

	resp.Result = result
	return nil
}
