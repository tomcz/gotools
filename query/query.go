package query

import (
	"context"
	"database/sql"
)

type ScanFn func(dest ...interface{}) error

type EachRowFn func(row ScanFn) error

type PartialQuery func(row EachRowFn) error

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
