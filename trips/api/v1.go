package api

import (
	"net/http"

	"starling/trips"

	"github.com/labstack/echo/v4"
)

type TripsAPIHandler struct {
	tripService *trips.TripService
}

func (t *TripsAPIHandler) CreateTrip(c echo.Context) error {
	// TODO: Validate dates etc.
	data := new(trips.TripData)
	if err := c.Bind(data); err != nil {
		return err
	}

	trip, err := t.tripService.CreateTrip(data)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, trip)
}
