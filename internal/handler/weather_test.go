package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWeatherHandler_CityLookup(t *testing.T) {
	target := fmt.Sprintf("/geoapi/city/lookup/?location=%f,%f", 120.13026, 30.25961)
	r := httptest.NewRequest(http.MethodGet, target, nil)
	w := httptest.NewRecorder()
	NewWeatherHandler().CityLookup(w, r)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	t.Logf("%+v", string(body))
}
