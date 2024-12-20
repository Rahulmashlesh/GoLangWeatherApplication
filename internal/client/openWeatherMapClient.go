package client

import (
	"fmt"
	"net/http"
)

type HttpGetter interface {
	Get(zipcode string, unit string) (*http.Response, error)
}

type OpenWeatherMapClient struct {
	apiKey string
}

func NewOpenWeatherMapClient(apikey string) *OpenWeatherMapClient {
	return &OpenWeatherMapClient{
		apiKey: apikey,
	}
}

func (c *OpenWeatherMapClient) Get(zipcode string, unit string) (rsp *http.Response, err error) {
	currentWeatherByZipUrl := "https://api.openweathermap.org/data/2.5/weather?zip=%s&appid=%s&units=%s"
	url := fmt.Sprintf(currentWeatherByZipUrl, zipcode, c.apiKey, unit)
	return http.Get(url)
}
