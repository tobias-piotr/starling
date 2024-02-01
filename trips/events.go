package trips

type TripEvent int64

const (
	TripCreated TripEvent = iota
	TripRequested
)

func (e TripEvent) String() string {
	return [...]string{"trip_created", "trip_requested"}[e]
}
