-- +goose Up
-- +goose StatementBegin
CREATE TABLE trip_results (
  id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  trip_id UUID NOT NULL REFERENCES trips(id) ON DELETE CASCADE,
  summary TEXT,
  attractions TEXT,
  weather TEXT,
  prices TEXT,
  luggage TEXT,
  documents TEXT,
  commuting TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE trip_results;
-- +goose StatementEnd
