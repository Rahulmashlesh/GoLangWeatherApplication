package weather

import (
	"GoWeaterAPI/internal/client"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

/*
var apiKey = "c4df6fe3e0d87fb4e9d14412929130a5" //   "os.Getenv("OWM_API_KEY")
var unit = "F"
var lang =
var currentWeatherByZipUrl = "https://api.openweathermap.org/data/2.5/weather?zip=%s&appid=%s"
*/
type CurrentWeather struct {
	Coord      Coord     `json:"coord"`
	Weather    []Weather `json:"weather"`
	Base       string    `json:"base"`
	Main       Main      `json:"main"`
	Visibility int       `json:"visibility"`
	Wind       Wind      `json:"wind"`
	Rain       Rain      `json:"rain"`
	Clouds     Clouds    `json:"clouds"`
	Dt         int       `json:"dt"`
	Sys        Sys       `json:"sys"`
	Timezone   int       `json:"timezone"`
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Cod        int       `json:"cod"`
	Client     client.HttpGetter
	Logger     *slog.Logger
	Zipcode    string
}
type Coord struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}
type Weather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}
type Main struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  int     `json:"pressure"`
	Humidity  int     `json:"humidity"`
	SeaLevel  int     `json:"sea_level"`
	GrndLevel int     `json:"grnd_level"`
}
type Wind struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
	Gust  float64 `json:"gust"`
}
type Rain struct {
	OneH float64 `json:"1h"`
}
type Clouds struct {
	All int `json:"all"`
}
type Sys struct {
	Type    int    `json:"type"`
	ID      int    `json:"id"`
	Country string `json:"country"`
	Sunrise int    `json:"sunrise"`
	Sunset  int    `json:"sunset"`
}

func NewCurrentWeather(httpGetter client.HttpGetter, logger *slog.Logger, zipcode string) *CurrentWeather {
	return &CurrentWeather{
		Zipcode: zipcode,
		Client:  httpGetter,
		Logger:  logger.With("context", "currentWeather", "zipcode", zipcode),
	}
}

func (w *CurrentWeather) Call() {
	w.GetWeather()
}

func (w *CurrentWeather) GetWeather() error {
	rsp, err := w.Client.Get(w.Zipcode)
	if err != nil {
		w.Logger.Error("Error during HTTP GET req:", err)
		return errors.New("Http Client Error")
	}

	w.Logger.Info("Response Code", "Rsp Status:", rsp.Status)

	if rsp.StatusCode == http.StatusOK {
		err = json.NewDecoder(rsp.Body).Decode(w)
		if err != nil {
			w.Logger.Error("Error Decoding Json", err)
			return errors.New("json Decoding Error")
		}
		w.Logger.Info("Received Weather", "Weather", w)
		return nil
	} else {
		// Handle non-OK status codes
		w.Logger.Error("Non-OK HTTP status code", "StatusCode", rsp.StatusCode)
		return errors.New("Http Client Error") // Change this line to return "Http Client Error"
	}
}
