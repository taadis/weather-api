package handler

import (
	"context"
	"testing"

	"github.com/taadis/weather-api/internal/model"
)

var (
	ctx     = context.Background()
	weather = NewWeather()
)

// 测试热门城市
func TestWeather_TopCity(t *testing.T) {
	req := &model.TopCityRequest{}
	resp := &model.TopCityResponse{}
	err := weather.TopCity(ctx, req, resp)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", resp)
}

// 测试城市信息查询
func TestWeather_LookupCity(t *testing.T) {
	req := &model.LookupCityRequest{
		Location: "120.13026,30.25961",
	}
	resp := &model.LookupCityResponse{}
	err := weather.LookupCity(ctx, req, resp)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", resp)
}

// 测试天气预报
func TestWeather_Forecast(t *testing.T) {
	req := &model.WeatherForecastRequest{
		Location: "120.13026,30.25961",
		Duration: "7d",
	}
	resp := &model.WeatherForecastResponse{}
	err := weather.Forecast(ctx, req, resp)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", resp)
}

// 测试实时天气
func TestWeatherHandler_Now(t *testing.T) {
	req := &model.WeatherNowRequest{
		Location: "120.13026,30.25961",
	}
	resp := &model.WeatherNowResponse{}
	err := weather.Now(ctx, req, resp)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", resp)
}

// 测试天气生活指数
func TestWeather_Indices(t *testing.T) {
	req := &model.WeatherIndicesRequest{
		Location: "120.13026,30.25961",
		Type:     "3,5",
		Duration: "1d",
	}
	resp := &model.WeatherIndicesResponse{}
	err := weather.Indices(ctx, req, resp)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", resp)
}
