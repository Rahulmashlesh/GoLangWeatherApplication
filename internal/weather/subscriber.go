package weather

type Subscriber interface {
	Notify(weather *CurrentWeather)
}
