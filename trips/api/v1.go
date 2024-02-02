package api

import (
	"net/http"

	"starling/trips"

	"github.com/labstack/echo/v4"
)

type TripsAPIHandler struct{}

func (t *TripsAPIHandler) CreateTrip(c echo.Context) error {
	// TODO: Validate dates etc.
	data := new(trips.TripData)
	if err := c.Bind(data); err != nil {
		return err
	}

	// TODO: Inject this
	srv := trips.TripService{}

	trip, err := srv.CreateTrip(data)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, trip)
}
