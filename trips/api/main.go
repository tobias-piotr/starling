package api

import "github.com/labstack/echo/v4"

func Register(g *echo.Group) {
	v1 := g.Group("/api/trips/v1")

	t := &TripsAPIHandler{}
	v1.POST("/", t.CreateTrip)
}
