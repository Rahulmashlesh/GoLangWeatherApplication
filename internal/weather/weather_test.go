package weather_test

import (
	"GoWeaterAPI/internal/weather"
	"github.com/stretchr/testify/assert"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"testing"
)

var currentWeatherByZipUrl = ""

func TestCurrentWeather_Call(t *testing.T) {
	zipcode := "95134"
	loggerForTest := *slog.With("context", "currentWeather", "zipcode", zipcode)
	client := httpMock{statusCode: 200, responseMessage: sample_rsp}
	weather := weather.NewCurrentWeather(&client, &loggerForTest, zipcode)
	err := weather.GetWeather()

	assert.NoError(t, err)
	assert.Equal(t, "Test Land", weather.Name)
}

func TestInvalidAPIKey(t *testing.T) {
	zipcode := "00000"
	loggerForTest := *slog.With("context", "currentWeather", "zipcode", zipcode)
	client := &httpMock{statusCode: 401, responseMessage: bad_apikey_rsp} // Indicate that the API key is invalid
	weather := weather.NewCurrentWeather(client, &loggerForTest, zipcode)
	err := weather.GetWeather()

	assert.Error(t, err)

}

func TestHTTPClientError(t *testing.T) {
	zipcode := "95134"
	loggerForTest := *slog.With("context", "currentWeather", "zipcode", zipcode)

	client := &httpMock{statusCode: http.StatusInternalServerError, responseMessage: ""}
	weatherObj := weather.NewCurrentWeather(client, &loggerForTest, zipcode)

	err := weatherObj.GetWeather()

	assert.Error(t, err)
	assert.Equal(t, "Http Client Error", err.Error())
}

type httpMock struct {
	statusCode      int
	responseMessage string
}

func (c *httpMock) Get(zipcode string) (rsp *http.Response, err error) {
	b := strings.NewReader(c.responseMessage)
	r := io.NopCloser(b)

	a := &http.Response{Body: r, StatusCode: c.statusCode, Status: http.StatusText(c.statusCode)}
	return a, nil
}

var bad_apikey_rsp = "Invalid API key"
var sample_rsp string = `
{
  "coord": {
    "lon": 10.99,
    "lat": 44.34
  },
  "weather": [
    {
      "id": 501,
      "main": "Rain",
      "description": "moderate rain",
      "icon": "10d"
    }
  ],
  "base": "stations",
  "main": {
    "temp": 298.48,
    "feels_like": 298.74,
    "temp_min": 297.56,
    "temp_max": 300.05,
    "pressure": 1015,
    "humidity": 64,
    "sea_level": 1015,
    "grnd_level": 933
  },
  "visibility": 10000,
  "wind": {
    "speed": 0.62,
    "deg": 349,
    "gust": 1.18
  },
  "rain": {
    "1h": 3.16
  },
  "clouds": {
    "all": 100
  },
  "dt": 1661870592,
  "sys": {
    "type": 2,
    "id": 2075663,
    "country": "IT",
    "sunrise": 1661834187,
    "sunset": 1661882248
  },
  "timezone": 7200,
  "id": 3163858,
  "name": "Test Land",
  "cod": 200
}               
`
