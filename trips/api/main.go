package api

import (
	"os"

	"starling/internal/events"
	"starling/trips"
	"starling/trips/database"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

func NewTripsAPIHandler(db *sqlx.DB, redisClient *redis.Client) *TripsAPIHandler {
	tripRepo := database.NewTripRepository(db)
	eventBus := events.NewRedisEventBus(
		redisClient,
		&events.RedisBusArgs{Stream: os.Getenv("REDIS_STREAM")},
	)
	tripService := trips.NewTripService(tripRepo, eventBus)
	return &TripsAPIHandler{tripService: tripService}
}

func Register(g *echo.Group, db *sqlx.DB, redisClient *redis.Client) {
	v1 := g.Group("/api/v1/trips")

	h := NewTripsAPIHandler(db, redisClient)
	v1.POST("", h.CreateTrip)
	v1.GET("", h.GetTrips)
	v1.GET("/:id", h.GetTrip)
}
