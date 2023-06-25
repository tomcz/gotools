package txx

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

// InTxx starts a context-aware database transaction wrapper, executes the
// callback function, and either commits the transaction if the callback
// exits without an error, or rolls-back the transaction if the callback
// returns an error.
func InTxx(ctx context.Context, db *sqlx.DB, callback func(tx *sqlx.Tx) error, opts ...*sql.TxOptions) error {
	var txOpts *sql.TxOptions
	for _, opt := range opts {
		if txOpts == nil {
			txOpts = opt
		} else {
			txOpts.ReadOnly = opt.ReadOnly
			txOpts.Isolation = opt.Isolation
		}
	}
	tx, err := db.BeginTxx(ctx, txOpts)
	if err != nil {
		return err
	}
	err = callback(tx)
	if err != nil {
		if ex := tx.Rollback(); ex != nil {
			return errors.Join(err, ex)
		}
		return err
	}
	return tx.Commit()
}
