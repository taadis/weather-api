package main

import (
	"os"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/etcd"
	"github.com/micro/go-micro/util/log"
	weatherSdk "github.com/taadis/qweather-sdk-go"
	"github.com/taadis/weather-api/internal/cache"
	"github.com/taadis/weather-api/internal/handler"
)

func main() {
	redisServer, err := miniredis.Run()
	if err != nil {
		log.Errorf("start redis server error:%+v", err)
		panic(err)
	}
	log.Infof("start redis server success")

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisServer.Addr(),
		Password: "",
		DB:       0,
	})

	iCache := cache.NewCache(rdb)
	iWeatherCache := handler.NewWeatherCache(iCache, weatherSdk.NewClient())

	service := micro.NewService(
		micro.BeforeStart(func() error {
			return nil
		}),
		micro.Name("go.micro.api.weather"),
		micro.Registry(etcd.NewRegistry(registry.Addrs(os.Getenv("MICRO_REGISTRY_ADDRESS")))),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
	)

	// optionally setup command line usage
	service.Init()

	// Register Handlers
	//proto.RegisterWeatherHandler(service.Server(), handler.NewWeatherHandler())
	micro.RegisterHandler(service.Server(), handler.NewWeather(iWeatherCache))

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
