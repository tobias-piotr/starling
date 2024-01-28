-- +goose Up
-- +goose StatementBegin
CREATE TABLE trips (
	id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	destination TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE trips;
-- +goose StatementEnd
