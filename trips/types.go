package trips

import (
	"fmt"
	"strings"
	"time"
)

// TripStatus is an enum representing the status of a trip
type TripStatus int64

const (
	Draft TripStatus = iota
	Requested
	Failed
	Completed
)

func (s TripStatus) String() string {
	return [...]string{"draft", "requested", "failed", "completed"}[s]
}

func (s TripStatus) MarshalJSON() ([]byte, error) {
	status := fmt.Sprintf(`"%s"`, s.String())
	return []byte(status), nil
}

// Date is a custom type for time.Time, that gets limited to YYYY-MM-DD format
type Date struct{ time.Time }

func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")

	date, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}

	d.Time = date
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	date := d.Format("2006-01-02")
	date = fmt.Sprintf(`"%s"`, date)
	return []byte(date), nil
}
