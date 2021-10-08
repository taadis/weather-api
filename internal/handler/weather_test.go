package handler

import (
	"context"
	"testing"

	"github.com/taadis/weather-api/proto"
)

// 测试天气预报
func TestWeatherHandler_Forecast(t *testing.T) {
	ctx := context.Background()
	req := &proto.ForecastRequest{
		Location: "120.13026,30.25961",
		Duration: "7d",
	}
	resp := &proto.ForecastResponse{}
	weatherHandler := NewWeatherHandler()
	err := weatherHandler.Forecast(ctx, req, resp)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", resp.Result)
}

// 测试实时天气
func TestWeatherHandler_Now(t *testing.T) {
	ctx := context.Background()
	req := &proto.NowRequest{
		Location: "120.13026,30.25961",
	}
	resp := &proto.NowResponse{}
	weatherHandler := NewWeatherHandler()
	err := weatherHandler.Now(ctx, req, resp)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", resp.Result)
}

// 测试天气生活指数
func TestWeatherHandler_Indices(t *testing.T) {
	ctx := context.Background()
	req := &proto.IndicesRequest{
		Location: "120.13026,30.25961",
		Type:     "3,5",
		Duration: "1d",
	}
	resp := &proto.IndicesResponse{}
	weatherHandler := NewWeatherHandler()
	err := weatherHandler.Indices(ctx, req, resp)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", resp.Result)
}
