package database

import (
	"starling/trips"

	"github.com/jmoiron/sqlx"
)

type TripRepository struct {
	db *sqlx.DB
}

func NewTripRepository(db *sqlx.DB) *TripRepository {
	return &TripRepository{db}
}

func (r *TripRepository) Create(data *trips.TripData) (*trips.Trip, error) {
	return &trips.Trip{}, nil
}
