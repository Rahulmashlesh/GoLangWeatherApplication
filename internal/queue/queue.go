package queue

import "context"

type Queue interface {
	Enqueue(context.Context, string) error
	Dequeue(context.Context) (string, error)
	Next(context.Context) (string, error)
	Delete(context.Context, string) error
	//IsEmpty(context.Context) (string, bool)
	//Size(context.Context) int
}
