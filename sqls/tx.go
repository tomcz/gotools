package sqls

import (
	"context"
	"database/sql"
	"errors"
)

// InTx starts a database transaction, executes the callback function,
// and either commits the transaction if the callback exits without an
// error, or rolls-back the transaction if the callback returns an error.
func InTx(db *sql.DB, callback func(tx *sql.Tx) error) error {
	tx, err := db.Begin()
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

// InTxContext starts a context-aware database transaction, executes the
// callback function, and either commits the transaction if the callback
// exits without an error, or rolls-back the transaction if the callback
// returns an error.
func InTxContext(ctx context.Context, db *sql.DB, callback func(tx *sql.Tx) error, opts ...*sql.TxOptions) error {
	var txOpts *sql.TxOptions
	for _, opt := range opts {
		txOpts = opt
	}
	tx, err := db.BeginTx(ctx, txOpts)
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
