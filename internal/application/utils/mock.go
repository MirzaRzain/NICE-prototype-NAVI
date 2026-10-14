package utils

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// NewMockDB creates a new mock database connection for testing purposes using sqlmock.
// It returns a *sql.DB instance, a sqlmock.Sqlmock object for setting expectations,
// and a cleanup function that should be deferred to ensure all expectations are met and the database is closed.
// This function is useful for unit testing code that interacts with a SQL database without requiring a real database connection.
func NewMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	cleanup := func() {
		assert.NoError(t, mock.ExpectationsWereMet())
		_ = db.Close()
	}

	return db, mock, cleanup
}
