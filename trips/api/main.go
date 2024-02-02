package api

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func Register(g *echo.Group, _ *sqlx.DB) {
	v1 := g.Group("/api/trips/v1")

	t := &TripsAPIHandler{}
	v1.POST("/", t.CreateTrip)
}
