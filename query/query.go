package query

import (
	"context"
	"database/sql"
)

// ScanFn provides an interface for the database/sql/#Rows.Scan
// function so that we can limit what is exposed by EachRowFn.
type ScanFn func(dest ...interface{}) error

// EachRowFn is called to process each query result row.
type EachRowFn func(row ScanFn) error

// PartialQuery is a curried query that passes each result row to EachRowFn.
type PartialQuery func(row EachRowFn) error

// QueryRows provides the entry point to retrieve a number of rows from
// a given query and arguments, using database/sql/#DB.QueryRow as inspiration.
func QueryRows(db *sql.DB, query string, args ...interface{}) PartialQuery {
	return func(row EachRowFn) error {
		rows, err := db.Query(query, args...)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			if err := row(rows.Scan); err != nil {
				return err
			}
		}
		return rows.Err()
	}
}

// QueryRowsContext provides the entry point to retrieve a number of rows from
// a given query and arguments, using database/sql/#DB.QueryRowContext as inspiration.
func QueryRowsContext(ctx context.Context, db *sql.DB, query string, args ...interface{}) PartialQuery {
	return func(row EachRowFn) error {
		rows, err := db.QueryContext(ctx, query, args...)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			if err := row(rows.Scan); err != nil {
				return err
			}
		}
		return rows.Err()
	}
}
