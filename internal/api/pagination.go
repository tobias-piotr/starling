package api

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

// GetPagination is a helper function to get the pagination parameters from the request.
func GetPagination(c echo.Context) (int, int, error) {
	page := c.QueryParam("page")
	if page == "" {
		page = "1"
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return 0, 0, err
	}
	if pageInt < 1 {
		pageInt = 1
	}

	perPage := c.QueryParam("per_page")
	if perPage == "" {
		perPage = "10"
	}
	perPageInt, err := strconv.Atoi(perPage)
	if err != nil {
		return 0, 0, err
	}
	if perPageInt < 1 {
		perPageInt = 10
	}

	return pageInt, perPageInt, nil
}
