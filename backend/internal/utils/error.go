package utils

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// For recognizing errors returned from PostgreSQL
func AssertPgErr(err error) (*pgconn.PgError, bool) {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		return pgErr, true
	}

	return nil, false
}
