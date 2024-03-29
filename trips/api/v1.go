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

// @Summary Create trip
// @Description Create a new trip
// @Tags trips
// @Param data body trips.TripData true "Request body"
// @Success 201 {object} trips.Trip
// @Router /api/v1/trips [post]
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

// @Summary Get trips
// @Description Get a list of trips
// @Tags trips
// @Success 200 {object} []trips.TripOverview
// @Router /api/v1/trips [get]
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

// @Summary Get trip
// @Description Get trip by id
// @Tags trips
// @Success 200 {object} []trips.Trip
// @Param id path string true "Trip ID"
// @Router /api/v1/trips/{id} [get]
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

// @Summary Request trip
// @Description Request trip result
// @Tags trips
// @Success 204
// @Param id path string true "Trip ID"
// @Router /api/v1/trips/{id}/request [post]
func (t *TripsAPIHandler) RequestTrip(c echo.Context) error {
	id := c.Param("id")

	if err := t.tripService.RequestTrip(id); err != nil {
		return api.RespondWithError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}
