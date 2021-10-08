// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/weather.proto

package proto

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	math "math"
)

import (
	context "context"
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for Weather service

type WeatherService interface {
	// 天气生活指数
	Indices(ctx context.Context, in *IndicesRequest, opts ...client.CallOption) (*IndicesResponse, error)
	// 实时天气
	Now(ctx context.Context, in *NowRequest, opts ...client.CallOption) (*NowResponse, error)
	// 逐天天气预报
	Forecast(ctx context.Context, in *ForecastRequest, opts ...client.CallOption) (*ForecastResponse, error)
}

type weatherService struct {
	c    client.Client
	name string
}

func NewWeatherService(name string, c client.Client) WeatherService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "proto"
	}
	return &weatherService{
		c:    c,
		name: name,
	}
}

func (c *weatherService) Indices(ctx context.Context, in *IndicesRequest, opts ...client.CallOption) (*IndicesResponse, error) {
	req := c.c.NewRequest(c.name, "Weather.Indices", in)
	out := new(IndicesResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *weatherService) Now(ctx context.Context, in *NowRequest, opts ...client.CallOption) (*NowResponse, error) {
	req := c.c.NewRequest(c.name, "Weather.Now", in)
	out := new(NowResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *weatherService) Forecast(ctx context.Context, in *ForecastRequest, opts ...client.CallOption) (*ForecastResponse, error) {
	req := c.c.NewRequest(c.name, "Weather.Forecast", in)
	out := new(ForecastResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Weather service

type WeatherHandler interface {
	// 天气生活指数
	Indices(context.Context, *IndicesRequest, *IndicesResponse) error
	// 实时天气
	Now(context.Context, *NowRequest, *NowResponse) error
	// 逐天天气预报
	Forecast(context.Context, *ForecastRequest, *ForecastResponse) error
}

func RegisterWeatherHandler(s server.Server, hdlr WeatherHandler, opts ...server.HandlerOption) error {
	type weather interface {
		Indices(ctx context.Context, in *IndicesRequest, out *IndicesResponse) error
		Now(ctx context.Context, in *NowRequest, out *NowResponse) error
		Forecast(ctx context.Context, in *ForecastRequest, out *ForecastResponse) error
	}
	type Weather struct {
		weather
	}
	h := &weatherHandler{hdlr}
	return s.Handle(s.NewHandler(&Weather{h}, opts...))
}

type weatherHandler struct {
	WeatherHandler
}

func (h *weatherHandler) Indices(ctx context.Context, in *IndicesRequest, out *IndicesResponse) error {
	return h.WeatherHandler.Indices(ctx, in, out)
}

func (h *weatherHandler) Now(ctx context.Context, in *NowRequest, out *NowResponse) error {
	return h.WeatherHandler.Now(ctx, in, out)
}

func (h *weatherHandler) Forecast(ctx context.Context, in *ForecastRequest, out *ForecastResponse) error {
	return h.WeatherHandler.Forecast(ctx, in, out)
}
