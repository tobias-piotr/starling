package tests

import (
	"os"
	"strings"

	"starling/internal/database"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
)

// SuiteWithDB is a test suite that handles connection with the test database
type SuiteWithDB struct {
	suite.Suite
	DB *sqlx.DB
}

// SetupSuite connects to the test database and runs migrations
func (s *SuiteWithDB) SetupSuite() {
	dsn := strings.Replace(os.Getenv("DATABASE_DSN"), "dbname=starling", "dbname=test", 1)
	db, _ := database.GetConnection(dsn)
	s.DB = db
	database.Migrate(s.DB)
}

// TearDownTest drops the public schema, recreates it, and runs migrations
func (s *SuiteWithDB) TearDownTest() {
	q := `
		DROP SCHEMA public CASCADE;
		CREATE SCHEMA public;
	`
	s.DB.MustExec(q)
	database.Migrate(s.DB)
}
