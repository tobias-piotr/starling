package api

import (
	"starling/internal/events"
	"starling/trips"
	"starling/trips/database"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func NewTripsAPIHandler(db *sqlx.DB, redisPublisher *events.RedisEventBus) *TripsAPIHandler {
	tripRepo := database.NewTripRepository(db)
	tripService := trips.NewTripService(tripRepo, redisPublisher)
	return &TripsAPIHandler{tripService: tripService}
}

func Register(g *echo.Group, db *sqlx.DB, redisPublisher *events.RedisEventBus) {
	v1 := g.Group("/api/v1/trips")

	h := NewTripsAPIHandler(db, redisPublisher)
	v1.POST("", h.CreateTrip)
	v1.GET("", h.GetTrips)
	v1.GET("/:id", h.GetTrip)
	v1.POST("/:id/request", h.RequestTrip)
}
