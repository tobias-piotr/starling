package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"starling/internal/events"
	"starling/internal/tests"
	"starling/trips"

	"github.com/stretchr/testify/suite"
)

type TripsAPISuite struct {
	tests.SuiteWithDB
	EventPublisher *events.RedisEventBus
}

func (s *TripsAPISuite) SetupSuite() {
	s.SuiteWithDB.SetupSuite()
	s.EventPublisher = tests.GetRedisPublisher()
}

func TestTripsAPISuite(t *testing.T) {
	suite.Run(t, new(TripsAPISuite))
}

func (s *TripsAPISuite) TestCreateTrip() {
	data := `{
		"name": "Test Trip",
		"destination": "Norway",
		"origin": "Wroclaw",
		"date_from": "2021-01-01",
		"date_to": "2021-01-09",
		"budget": 1000,
		"requirements": "Cannot be uncomfy."
	}`

	rec, c := tests.MakeRequest(http.MethodPost, "/api/v1/trips", "", data)
	h := NewTripsAPIHandler(s.DB, s.EventPublisher)

	if s.NoError(h.CreateTrip(c)) {
		s.Equal(http.StatusCreated, rec.Code)
		trip := trips.Trip{}
		err := json.Unmarshal(rec.Body.Bytes(), &trip)
		s.NoError(err)
		s.NotEmpty(trip.ID)
	}
}

func (s *TripsAPISuite) TestCreateTripValidation() {
	data := `{
		"name": "",
		"destination": "",
		"origin": "",
		"date_from": "2021-01-02",
		"date_to": "2021-01-01",
		"budget": -1
	}`

	rec, c := tests.MakeRequest(http.MethodPost, "/api/v1/trips", "", data)
	h := NewTripsAPIHandler(s.DB, s.EventPublisher)

	if s.NoError(h.CreateTrip(c)) {
		s.Equal(http.StatusBadRequest, rec.Code)
		errs := map[string][]string{}
		err := json.Unmarshal(rec.Body.Bytes(), &errs)
		s.NoError(err)
		s.Equal([]string{"Name can't be blank"}, errs["name"])
		s.Equal([]string{"Destination can't be blank"}, errs["destination"])
		s.Equal([]string{"Origin can't be blank"}, errs["origin"])
		s.Equal([]string{"Date from must be before date to"}, errs["date_from"])
		s.Equal([]string{`Budget must be greater than "0"`}, errs["budget"])

	}
}

func (s *TripsAPISuite) TestGetTrips() {
	data := `{
		"name": "Test Trip",
		"destination": "Norway",
		"origin": "Wroclaw",
		"date_from": "2021-01-01",
		"date_to": "2021-01-09",
		"budget": 1000,
		"requirements": "Cannot be uncomfy."
	}`

	h := NewTripsAPIHandler(s.DB, s.EventPublisher)

	// Create 5 trips
	for i := 1; i <= 5; i++ {
		data = strings.Replace(data, "Test Trip", fmt.Sprintf("Test Trip %v", i), 1)
		_, c := tests.MakeRequest(http.MethodPost, "/api/v1/trips", "", data)
		s.NoError(h.CreateTrip(c))
	}

	// Get all trips
	rec, c := tests.MakeRequest(http.MethodGet, "/api/v1/trips", "", "")
	if s.NoError(h.GetTrips(c)) {
		s.Equal(http.StatusOK, rec.Code)
		trips := []trips.Trip{}
		err := json.Unmarshal(rec.Body.Bytes(), &trips)
		s.NoError(err)
		s.Len(trips, 5)
	}
}

func (s *TripsAPISuite) TestGetTrip() {
	data := `{
		"name": "Test Trip",
		"destination": "Norway",
		"origin": "Wroclaw",
		"date_from": "2021-01-01",
		"date_to": "2021-01-09",
		"budget": 1000,
		"requirements": "Cannot be uncomfy."
	}`

	h := NewTripsAPIHandler(s.DB, s.EventPublisher)

	// Create a trip
	rec, c := tests.MakeRequest(http.MethodPost, "/api/v1/trips", "", data)
	s.NoError(h.CreateTrip(c))
	createdTrip := trips.Trip{}
	err := json.Unmarshal(rec.Body.Bytes(), &createdTrip)
	s.NoError(err)

	// Get the trip
	rec, c = tests.MakeRequest(http.MethodGet, "/api/v1/trips", createdTrip.ID.String(), "")
	if s.NoError(h.GetTrip(c)) {
		s.Equal(http.StatusOK, rec.Code)
		trip := trips.Trip{}
		err := json.Unmarshal(rec.Body.Bytes(), &trip)
		s.NoError(err)
		s.NotEmpty(trip.ID)
	}
}
