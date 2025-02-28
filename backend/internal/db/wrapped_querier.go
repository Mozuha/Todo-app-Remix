package db

import (
	"github.com/jackc/pgx/v5"
)

// https://github.com/sqlc-dev/sqlc/issues/383

type WrappedQuerier interface {
	Querier
	WithTx(tx pgx.Tx) WrappedQuerier
}

type WrappedQueries struct {
	*Queries
}

func (q *WrappedQueries) WithTx(tx pgx.Tx) WrappedQuerier {
	return &WrappedQueries{
		Queries: q.Queries.WithTx(tx),
	}
}

func NewWrappedQuerier(q *Queries) WrappedQuerier {
	return &WrappedQueries{
		Queries: q,
	}
}
