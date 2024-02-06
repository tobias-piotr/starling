-- +goose Up
-- +goose StatementBegin
CREATE TYPE trip_status AS ENUM ('draft', 'requested', 'failed', 'completed');
CREATE TABLE trips (
	id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	status trip_status DEFAULT 'draft',
	name VARCHAR(255) NOT NULL,
	destination TEXT NOT NULL,
	origin TEXT NOT NULL,
	date_from DATE,
	date_to DATE,
	budget INT,
	requirements TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE trips;
-- +goose StatementEnd
