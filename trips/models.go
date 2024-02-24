package trips

import (
	"time"

	"github.com/cohesivestack/valgo"
	"github.com/google/uuid"
)

type Trip struct {
	ID           uuid.UUID   `json:"id"`
	CreatedAt    time.Time   `json:"created_at" db:"created_at"`
	Status       TripStatus  `json:"status"`
	Name         string      `json:"name"`
	Destination  string      `json:"destination"`
	Origin       string      `json:"origin"`
	DateFrom     Date        `json:"date_from" db:"date_from"`
	DateTo       Date        `json:"date_to" db:"date_to"`
	Budget       int64       `json:"budget"`
	Requirements string      `json:"requirements"`
	Result       *TripResult `json:"result"`
}

// ValidateRequest validates the state of the trip before it can be requested.
// Validation skips the fields that are enforced elsewhere, like id or name.
func (t *Trip) ValidateRequest() error {
	if err := valgo.Is(valgo.Int64(t.Status, "status").EqualTo(DraftStatus)).Error(); err != nil {
		return err
	}

	return valgo.
		Is(valgo.String(t.Destination, "destination").Not().Blank()).
		Is(valgo.String(t.Origin, "origin").Not().Blank()).
		Is(valgo.Int64(t.Budget, "budget").GreaterThan(0)).
		Is(valgo.Any(t.DateFrom.NullableString(), "date_from").Not().Nil()).
		Is(valgo.Any(t.DateTo.NullableString(), "date_to").Not().Nil()).
		Is(valgo.String(t.Requirements, "requirements").Not().Blank()).
		Error()
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
	v := valgo.Is(valgo.String(d.Name, "name").Not().Blank())

	if d.DateFrom.After(d.DateTo.Time) {
		v.AddErrorMessage("date_from", "Date from must be before date to")
	}

	return v.Error()
}

type TripResult struct {
	ID          uuid.UUID `json:"id"`
	Summary     string    `json:"summary"`
	Attractions string    `json:"attractions"`
	Weather     string    `json:"weather"`
	Prices      string    `json:"prices"`
	Luggage     string    `json:"luggage"`
	Documents   string    `json:"documents"`
	Commuting   string    `json:"commuting"`
}
