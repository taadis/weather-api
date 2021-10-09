package model

// WeatherNowRequest 实时天气请求
type WeatherNowRequest struct {
	Location string `json:"location,omitempty"`
}

// WeatherNowResponse 实时天气响应
type WeatherNowResponse map[string]interface{}
