package database

import "starling/trips"

type TripRepository interface {
	Create(data *trips.TripData) (*trips.Trip, error)
}
