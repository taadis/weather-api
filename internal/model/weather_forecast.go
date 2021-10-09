package model

// WeatherForecastRequest 逐天天气预报请求
type WeatherForecastRequest struct {
	Location string `json:"location,omitempty"`
	Duration string `json:"duration,omitempty"`
}

// WeatherForecastResponse 逐天天气预报响应
type WeatherForecastResponse map[string]interface{}
