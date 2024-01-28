package database

import (
	"embed"

	"github.com/jmoiron/sqlx"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var fs embed.FS

// Migrate uses goose to read the migrations from the embedded filesystem and
// apply them to the database.
func Migrate(conn *sqlx.DB) error {
	goose.SetBaseFS(fs)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Up(conn.DB, "migrations")
}
