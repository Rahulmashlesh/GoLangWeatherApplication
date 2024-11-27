package models

type Location struct {
	Zipcode     string  `json:"zipcode"`
	Name        string  `json:"name"`
	Temperature float64 `json:"temperature"`
	// add client mentircs and logger, implement factory func with these to return locaiton with these.
}
func NewLocation(zipcode string) *Location {
	return &Location{
		Zipcode:     zipcode,
		Name:        "",
		Temperature: float64(0),
	}
}
