package handler

import (
	"log"
	"net/http"
	"os"

	qweather "github.com/Ink-33/go-heweather/v7"
)

const (
	QWEATHER_PUBLIC_ID = "QWEATHER_PUBLIC_ID"
	QWEATHER_KEY       = "QWEATHER_KEY"
	isBusiness         = false // 免费开发版为false，商业共享版与商业高性能版均为true
)

type WeatherHandler struct {
	qweatherCredential *qweather.Credential
}

func NewWeatherHandler() *WeatherHandler {
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

// Now 实时天气
func (h *WeatherHandler) Now(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	location := query.Get("location")
	client := qweather.NewRealTimeWeatherClient(location)
	resp, err := client.Run(h.qweatherCredential, nil)
	if err != nil {
		log.Printf("got now weather error: %v", err) //也可以自行进行错误处理
	}
	w.Write([]byte(resp))
}

// Forecast 天气预报
func (h *WeatherHandler) Forecast(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	location := query.Get("location")
	duration := query.Get("duration")
	client, err := qweather.NewWeatherForecastClient(location, duration)
	if err != nil {
		log.Printf("qweather.NewWeatherForecastClient error: %v, location: %s, duration: %s", err, location, duration)
	}
	resp, err := client.Run(h.qweatherCredential, nil)
	if err != nil {
		log.Printf("client.Run error: %v", err)
	}
	w.Write([]byte(resp))
}

func (h *WeatherHandler) CityTop(w http.ResponseWriter, r *http.Request) {
	//rawQuery := r.URL.RawQuery
	client := qweather.NewGeoTopCityClient()
	resp, err := client.Run(h.qweatherCredential, nil)
	if err != nil {
		log.Printf("client.Run error: %v", err)
	}

	w.Write([]byte(resp))
}

// CityLookup 城市搜索服务
func (h *WeatherHandler) CityLookup(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	location := query.Get("location")
	client := qweather.NewGeoCityClient(location)
	resp, err := client.Run(h.qweatherCredential, nil)
	if err != nil {
		log.Printf("got city lookup error: %v", err)
	}
	w.Write([]byte(resp))
}
