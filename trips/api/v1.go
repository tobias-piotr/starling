package api

import (
	"fmt"
	"net/http"
	"time"

	"starling/trips"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TripsAPIHandler struct{}

func (t *TripsAPIHandler) CreateTrip(c echo.Context) error {
	data := new(trips.TripData)
	if err := c.Bind(data); err != nil {
		return err
	}
	// TODO: Validate dates etc.

	trip := trips.Trip{
		ID:           uuid.New(),
		CreatedAt:    time.Now(),
		Status:       trips.Draft,
		Name:         data.Name,
		Destination:  data.Destination,
		Origin:       data.Origin,
		DateFrom:     data.DateFrom,
		DateTo:       data.DateTo,
		Budget:       data.Budget,
		Requirements: data.Requirements,
	}
	// TODO: Emit event
	fmt.Println(trips.TripCreated)

	return c.JSON(http.StatusCreated, trip)
}
