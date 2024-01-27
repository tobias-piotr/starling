package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func GetConnection(dsn string) (*sqlx.DB, error) {
	return sqlx.Connect("postgres", dsn)
}
