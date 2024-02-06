package api

import (
	"starling/internal/events"
	"starling/trips"
	"starling/trips/database"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func NewTripsAPIHandler(db *sqlx.DB) *TripsAPIHandler {
	tripRepo := database.NewTripRepository(db)
	eventBus := events.NewRedisEventBus()
	tripService := trips.NewTripService(tripRepo, eventBus)
	return &TripsAPIHandler{tripService: tripService}
}

func Register(g *echo.Group, db *sqlx.DB) {
	v1 := g.Group("/api/trips/v1")

	h := NewTripsAPIHandler(db)
	v1.POST("/", h.CreateTrip)
	v1.GET("/", h.GetTrips)
}
