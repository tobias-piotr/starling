package database

import (
	"database/sql"
	"errors"

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
	query := `
	INSERT INTO trips (status, name, destination, origin, date_from, date_to, budget, requirements)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id, created_at, status, name, destination, origin, date_from, date_to, budget, requirements;
	`
	var trip trips.Trip
	err := r.db.QueryRowx(
		query,
		trips.DraftStatus.String(),
		data.Name,
		data.Destination,
		data.Origin,
		data.DateFrom.Format("2006-01-02"),
		data.DateTo.Format("2006-01-02"),
		data.Budget,
		data.Requirements,
	).StructScan(&trip)
	if err != nil {
		return nil, err
	}

	return &trip, nil
}

func (r *TripRepository) GetAll(page int, perPage int) ([]*trips.TripOverview, error) {
	query := `
	SELECT id, created_at, status, name
	FROM trips
	ORDER BY created_at DESC
	LIMIT $1
	OFFSET $2;
	`
	offset := 0
	if page > 0 && offset > 0 {
		offset = (page - 1) * perPage
	}

	var trips []*trips.TripOverview
	err := r.db.Select(&trips, query, perPage, offset)
	if err != nil {
		return nil, err
	}

	return trips, nil
}

func (r *TripRepository) Get(id string) (*trips.Trip, error) {
	query := `
	SELECT id, created_at, status, name, destination, origin, date_from, date_to, budget, requirements
	FROM trips
	WHERE id = $1;
	`
	var trip trips.Trip
	err := r.db.QueryRowx(query, id).StructScan(&trip)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &trip, nil
}
