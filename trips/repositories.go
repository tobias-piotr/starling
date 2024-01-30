package trips

type TripRepository interface {
	Create(data *TripData) (*Trip, error)
}
