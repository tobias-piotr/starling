package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"starling/internal/database"
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
	if err := r.db.QueryRowx(
		query,
		trips.DraftStatus.String(),
		data.Name,
		data.Destination,
		data.Origin,
		data.DateFrom.NullableString(),
		data.DateTo.NullableString(),
		data.Budget,
		data.Requirements,
	).StructScan(&trip); err != nil {
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
	if err := r.db.Select(&trips, query, perPage, offset); err != nil {
		return nil, err
	}

	return trips, nil
}

func (r *TripRepository) Get(id string) (*trips.Trip, error) {
	query := `
	SELECT t.id, t.created_at, t.status, t.name, t.destination, t.origin, t.date_from, t.date_to, t.budget, t.requirements,
        CASE WHEN tr.id IS NOT NULL THEN
        json_build_object(
            'id', tr.id,
            'summary', tr.summary,
            'attractions', tr.attractions,
            'weather', tr.weather,
            'prices', tr.prices,
            'luggage', tr.luggage,
            'documents', tr.documents,
            'commuting', tr.commuting
        )
        ELSE NULL END as result
	FROM trips t
    LEFT JOIN trip_results tr ON t.id = tr.trip_id
	WHERE t.id = $1;
	`
	var trip trips.Trip
	var result sql.NullString
	err := r.db.QueryRowx(query, id).Scan(
		&trip.ID,
		&trip.CreatedAt,
		&trip.Status,
		&trip.Name,
		&trip.Destination,
		&trip.Origin,
		&trip.DateFrom,
		&trip.DateTo,
		&trip.Budget,
		&trip.Requirements,
		&result,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Unmarshal the result
	if result.String != "" {
		if err := json.Unmarshal([]byte(result.String), &trip.Result); err != nil {
			return nil, err
		}
	}
	return &trip, nil
}

func (r *TripRepository) Update(id string, data map[string]any) error {
	query := `
	UPDATE trips
	SET %s
	WHERE id = :id;
	`
	args := database.ConvertMapToArgsStr(data, ", ")
	query = fmt.Sprintf(query, args)
	data["id"] = id

	_, err := r.db.NamedExec(query, data)
	if err != nil {
		return err
	}

	return nil
}

func (r *TripRepository) AddResult(tripID string) (*trips.TripResult, error) {
	query := `
        INSERT INTO trip_results (trip_id)
        VALUES ($1)
        RETURNING id, trip_id;
        `

	var result trips.TripResult
	if err := r.db.QueryRowx(query, tripID).StructScan(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *TripRepository) UpdateResult(resultID string, data map[string]any) error {
	query := `
        UPDATE trip_results
        SET %s
        WHERE id = :id;
        `
	args := database.ConvertMapToArgsStr(data, ", ")
	query = fmt.Sprintf(query, args)
	data["id"] = resultID

	_, err := r.db.NamedExec(query, data)
	if err != nil {
		return err
	}

	return nil
}
