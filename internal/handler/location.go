package handlers

import (
	"GoWeatherAPI/internal/models"
	"GoWeatherAPI/internal/template"
	"context"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v5"
	"log/slog"
)

type Location struct {
	model  models.Model
	logger *slog.Logger
}

func (l *Location) List(c echo.Context) error {
	ctx := context.Background()

	_, err := l.model.List(ctx)
	if err != nil {
		l.logger.Error("Error getting all locations", "error", err)
		return err
	}

	component := template.Hello("Welcome Page", "rahul")
	h := templ.Handler(component)
	return h.Component.Render(ctx, c.Response().Writer)
}

func NewLocation(store *models.RedisStore, logger *slog.Logger) *Location {
	return &Location{
		model:  store,
		logger: logger,
	}
}
