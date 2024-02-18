package trips

import (
	"fmt"
	"strings"
	"time"
)

// TripStatus is an enum representing the status of a trip.
type TripStatus int64

const (
	DraftStatus TripStatus = iota
	RequestedStatus
	FailedStatus
	CompletedStatus
)

func (s TripStatus) Values() [4]string {
	return [...]string{"draft", "requested", "failed", "completed"}
}

// String returns the string representation of a TripStatus.
func (s TripStatus) String() string {
	return s.Values()[s]
}

// MarshalJSON is a custom marshaller to convert TripStatus to JSON.
func (s TripStatus) MarshalJSON() ([]byte, error) {
	status := fmt.Sprintf(`"%s"`, s.String())
	return []byte(status), nil
}

// UnmarshalJSON is a custom marshaller to convert a JSON string to TripStatus.
func (s *TripStatus) UnmarshalJSON(b []byte) error {
	v := strings.Trim(string(b), "\"")

	for i, status := range s.Values() {
		if status == v {
			*s = TripStatus(i)
			return nil
		}
	}

	return fmt.Errorf("unmarshal status: %v", v)
}

// Scan is used to convert a value from the database to a TripStatus.
func (s *TripStatus) Scan(value any) error {
	v, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("scan status: %v", value)
	}

	for i, status := range s.Values() {
		if status == string(v) {
			*s = TripStatus(i)
			return nil
		}
	}

	return fmt.Errorf("scan status: %v", value)
}

// Date is a custom type for time.Time, that gets limited to YYYY-MM-DD format.
type Date struct{ time.Time }

// NullableString returns a string representation of a Date that can be null.
func (d Date) NullableString() *string {
	if d.IsZero() {
		return nil
	}
	date := d.Format("2006-01-02")
	return &date
}

// UnmarshalJSON is a custom marshaller to convert a JSON string to a Date.
func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	date, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}

	d.Time = date
	return nil
}

// MarshalJSON is a custom marshaller to convert a Date to JSON.
func (d Date) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return []byte("null"), nil
	}
	date := d.Format("2006-01-02")
	date = fmt.Sprintf(`"%s"`, date)
	return []byte(date), nil
}

// Scan is used to convert a value from the database to a Date.
func (s *Date) Scan(value any) error {
	if value == nil {
		return nil
	}

	v, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("scan date: %v", value)
	}

	*s = Date{v}
	return nil
}
