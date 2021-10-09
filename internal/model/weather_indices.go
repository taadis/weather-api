package model

// WeatherIndicesRequest 天气生活指数请求
type WeatherIndicesRequest struct {
	Location string `json:"location,omitempty"`
	Type     string `json:"type,omitempty"`
	Duration string `json:"duration,omitempty"`
}

// WeatherIndicesResponse 天气生活指数响应
type WeatherIndicesResponse map[string]interface{}
