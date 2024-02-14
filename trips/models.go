package trips

import (
	"time"

	"github.com/cohesivestack/valgo"
	"github.com/google/uuid"
)

type Trip struct {
	ID           uuid.UUID  `json:"id"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	Status       TripStatus `json:"status"`
	Name         string     `json:"name"`
	Destination  string     `json:"destination"`
	Origin       string     `json:"origin"`
	DateFrom     Date       `json:"date_from" db:"date_from"`
	DateTo       Date       `json:"date_to" db:"date_to"`
	Budget       int64      `json:"budget"`
	Requirements string     `json:"requirements"`
}

type TripOverview struct {
	ID        uuid.UUID  `json:"id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	Status    TripStatus `json:"status"`
	Name      string     `json:"name"`
}

type TripData struct {
	Name         string `json:"name"`
	Destination  string `json:"destination"`
	Origin       string `json:"origin"`
	DateFrom     Date   `json:"date_from"`
	DateTo       Date   `json:"date_to"`
	Budget       int64  `json:"budget"`
	Requirements string `json:"requirements"`
}

func (d TripData) Validate() error {
	v := valgo.
		Is(valgo.String(d.Name, "name").Not().Blank()).
		Is(valgo.String(d.Destination, "destination").Not().Blank()).
		Is(valgo.String(d.Origin, "origin").Not().Blank()).
		Is(valgo.Int64(d.Budget, "budget").GreaterThan(0))

	if d.DateFrom.After(d.DateTo.Time) {
		v.AddErrorMessage("date_from", "Date from should be before date to")
	}

	return v.Error()
}

type TripResult struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	TripID      uuid.UUID
	Summary     string
	Attractions string
	Weather     string
	Prices      string
	Luggage     string
	Documents   string
	Commuting   string
}
