package sqls

import (
	"context"
	"database/sql"
)

// ScanFunc provides an interface for the database/sql/#Rows.Scan
// function so that we can limit what is exposed by EachRowFunc.
type ScanFunc func(dest ...interface{}) error

// EachRowFunc is called to process each query result row.
type EachRowFunc func(row ScanFunc) error

// PartialQuery is a curried query that passes each result row to EachRowFunc.
type PartialQuery func(row EachRowFunc) error

// QueryRows provides the entry point to retrieve a number of rows from
// a given query and arguments, using database/sql/#DB.QueryRow as inspiration.
func QueryRows(db *sql.DB, query string, args ...interface{}) PartialQuery {
	return func(row EachRowFunc) error {
		rows, err := db.Query(query, args...)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			err = row(rows.Scan)
			if err != nil {
				return err
			}
		}
		return rows.Err()
	}
}

// QueryRowsContext provides the entry point to retrieve a number of rows from
// a given query and arguments, using database/sql/#DB.QueryRowContext as inspiration.
func QueryRowsContext(ctx context.Context, db *sql.DB, query string, args ...interface{}) PartialQuery {
	return func(row EachRowFunc) error {
		rows, err := db.QueryContext(ctx, query, args...)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			err = row(rows.Scan)
			if err != nil {
				return err
			}
		}
		return rows.Err()
	}
}
