package weather_test

import (
	"GoWeaterAPI/internal/weather"
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

	client := httpMock{}
	// Create an instance of CurrentWeather with the fake logger and a dummy zipcode
	weather := weather.NewCurrentWeather(&client, &loggerForTest, zipcode)

	// Call the API (this will hit the fake server)
	weather.Call()

	// TODO assert and check the log content.

	//weather.Logger.

}

var sample_rsp = `
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

type httpMock struct {
}

func (c *httpMock) Get(zipcode string) (rsp *http.Response, err error) {

	r := strings.NewReader(sample_rsp)
	return &http.Response{Body: io.NopCloser(r), Status: http.StatusText(200)}, nil
}
