package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/taadis/weather-service/internal/handler"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Hello world received a request.")
	target := "Welcome to CloudBase"
	fmt.Fprintf(w, "Hello, %s!\n", target)
}

func main() {
	log.Print("Hello world sample started.")

	weatherHandler := handler.NewWeatherHandler()
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/api/weather/now/", weatherHandler.Now)
	http.HandleFunc("/api/weather/forecast/", weatherHandler.Forecast)
	http.HandleFunc("/geoapi/city/top/", weatherHandler.CityTop)
	http.HandleFunc("/geoapi/city/lookup/", weatherHandler.CityLookup)

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Println("http.ListenAndServe", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("http.ListenAndServe error: ", err)
	}
}
