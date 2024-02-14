package api

import (
	"net/http"

	"starling/internal/api"
	"starling/trips"

	"github.com/labstack/echo/v4"
)

type TripsAPIHandler struct {
	tripService *trips.TripService
}

func (t *TripsAPIHandler) CreateTrip(c echo.Context) error {
	data := new(trips.TripData)
	if err := c.Bind(data); err != nil {
		return err
	}
	if err := data.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	trip, err := t.tripService.CreateTrip(data)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, trip)
}

func (t *TripsAPIHandler) GetTrips(c echo.Context) error {
	page, perPage, err := api.GetPagination(c)
	if err != nil {
		return err
	}

	trips, err := t.tripService.GetTrips(page, perPage)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, trips)
}

func (t *TripsAPIHandler) GetTrip(c echo.Context) error {
	id := c.Param("id")

	trip, err := t.tripService.GetTrip(id)
	if err != nil {
		return err
	}

	if trip == nil {
		return c.JSON(http.StatusNotFound, nil)
	}

	return c.JSON(http.StatusOK, trip)
}
