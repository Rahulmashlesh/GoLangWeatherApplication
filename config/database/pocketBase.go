package database

import (
	"github.com/pocketbase/pocketbase"
)

var PocketBaseApp *pocketbase.PocketBase

func InitializeApp() {
	//PocketBaseApp = pocketbase.New()
}

/*
package models

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
)

type Location struct {
	Zipcode     string  `db:"Zipcode"`
	Name        string  `db:"Location"`
	Temperature float64 `db:"Temperature"`
}

type LocationStore interface {
	Get(string) (*Location, error)
	Create(*Location) error
	Update(*Location) error
	List(*Location) error
}

type Locations struct {
	dao *daos.Dao
}

func (l *Locations) Get(zipcode string) (*Location, error) {
	location := &Location{}
	err := l.dao.DB().
		Select("Zipcode", "Location", "Temperature").
		From("WeatherData2").
		Where(dbx.NewExp("Zipcode = {:zipcode}", dbx.Params{"Zipcode": zipcode})).
		One(location)
	return location, err
}

func (l *Locations) Create(location *Location) error {

	_, err := l.dao.DB().Insert("WeatherData2",
		dbx.Params{"Zipcode": location.Zipcode, "Location": location.Name, "Temperature": location.Temperature}).Execute()
	return err

}

func (l *Locations) Update(location *Location) error {

	updateParams := dbx.Params{
		"Location":    location.Name,
		"Temperature": location.Temperature,
	}
	whereClause := dbx.HashExp{"Zipcode": location.Zipcode}
	_, err := l.dao.DB().Update("WeatherData2", updateParams, whereClause).Execute()

	return err
}

func NewLocations(dao1 *daos.Dao) *Locations {
	a := &Locations{
		dao: dao1,
	}

	return a
}
*/
