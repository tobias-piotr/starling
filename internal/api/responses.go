package api

import (
	"net/http"

	"starling/internal/domain"

	"github.com/labstack/echo/v4"
)

// RespondWithError compares error type with known domain errors and returns the appropriate HTTP status code.
func RespondWithError(c echo.Context, err error) error {
	switch e := err.(type) {
	case *domain.NotFoundErr:
		return c.JSON(http.StatusNotFound, map[string]string{"message": e.Error()})
	case *domain.ValidationErr:
		return c.JSON(http.StatusBadRequest, e.Err)
	default:
		return err
	}
}
