package trips

type TripEvent int64

const (
	TripCreated TripStatus = iota
	TripRequested
)

func (s TripStatus) String() string {
	return [...]string{"trip_created", "trip_requested"}[s]
}
