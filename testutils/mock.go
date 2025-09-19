package testutils

import (
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

// MockDB creates a mock database for testing
func MockDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}

	sqlxDB := sqlx.NewDb(db, "postgres")
	return sqlxDB, mock
}

// MockLogger creates a test logger
func MockLogger(t *testing.T) *zap.Logger {
	return zaptest.NewLogger(t)
}

// MockRows creates mock rows for database queries
func MockRows(columns []string, values ...[]driver.Value) *sqlmock.Rows {
	rows := sqlmock.NewRows(columns)
	for _, value := range values {
		rows.AddRow(value...)
	}
	return rows
}

// MockError creates a mock database error
func MockError(err error) error {
	return err
}

// MockQuery creates a mock query expectation
func MockQuery(mock sqlmock.Sqlmock, query string, args []driver.Value, rows *sqlmock.Rows) {
	mock.ExpectQuery(query).WithArgs(args...).WillReturnRows(rows)
}

// MockExec creates a mock exec expectation
func MockExec(mock sqlmock.Sqlmock, query string, args []driver.Value, result driver.Result) {
	mock.ExpectExec(query).WithArgs(args...).WillReturnResult(result)
}

// MockBegin creates a mock transaction begin expectation
func MockBegin(mock sqlmock.Sqlmock) {
	mock.ExpectBegin()
}

// MockCommit creates a mock transaction commit expectation
func MockCommit(mock sqlmock.Sqlmock) {
	mock.ExpectCommit()
}

// MockRollback creates a mock transaction rollback expectation
func MockRollback(mock sqlmock.Sqlmock) {
	mock.ExpectRollback()
}
