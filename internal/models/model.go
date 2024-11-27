package models

import "context"

type Model interface {
	Get(context.Context, string) (*Location, error)
	Create(context.Context, *Location) error
	Update(context.Context, *Location) error
	List(context.Context) ([]Location, error)
	Delete(context.Context, string) error
}
