package handler

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/taadis/weather-api/proto"
)

func TestWeatherHandler_CityTop(t *testing.T) {
	ctx := context.Background()
	weatherHandler := NewWeatherHandler()
	req := &proto.TopCityRequest{}
	resp := &proto.TopCityResponse{}
	err := weatherHandler.TopCity(ctx, req, resp)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", resp)
}

func TestWeatherHandler_CityLookup(t *testing.T) {
	target := fmt.Sprintf("/geoapi/city/lookup/?location=%f,%f", 120.13026, 30.25961)
	r := httptest.NewRequest(http.MethodGet, target, nil)
	w := httptest.NewRecorder()
	NewWeatherHandler().CityLookup(w, r)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	t.Logf("%+v", string(body))
}
