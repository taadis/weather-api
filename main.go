package main

import (
	"log"
	"os"
	"time"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/etcd"
	"github.com/taadis/weather-api/internal/handler"
	"github.com/taadis/weather-api/proto"
)

func main() {
	service := micro.NewService(
		micro.Name("go.micro.api.weather"),
		micro.Registry(etcd.NewRegistry(registry.Addrs(os.Getenv("MICRO_REGISTRY_ADDRESS")))),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
	)

	// optionally setup command line usage
	service.Init()

	// Register Handlers
	proto.RegisterWeatherHandler(service.Server(), handler.NewWeatherHandler())

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

//
//import (
//	"fmt"
//	"log"
//	"net/http"
//	"os"
//
//	"github.com/taadis/weather-service/internal/handler"
//)
//
//func defaultHandler(w http.ResponseWriter, r *http.Request) {
//	log.Print("Hello world received a request.")
//	target := "Welcome to CloudBase"
//	fmt.Fprintf(w, "Hello, %s!\n", target)
//}
//
//func main() {
//	log.Print("Hello world sample started.")
//
//	weatherHandler := handler.NewWeatherHandler()
//	http.HandleFunc("/", defaultHandler)
//	http.HandleFunc("/api/weather/now/", weatherHandler.Now)
//	http.HandleFunc("/api/weather/forecast/", weatherHandler.Forecast)
//	http.HandleFunc("/geoapi/city/top/", weatherHandler.CityTop)
//	http.HandleFunc("/geoapi/city/lookup/", weatherHandler.CityLookup)
//
//	port := os.Getenv("PORT")
//	if port == "" {
//		port = "80"
//	}
//
//	addr := fmt.Sprintf(":%s", port)
//	log.Println("http.ListenAndServe", addr)
//	err := http.ListenAndServe(addr, nil)
//	if err != nil {
//		log.Fatal("http.ListenAndServe error: ", err)
//	}
//}
