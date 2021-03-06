package sqls

import (
	"context"
	"database/sql"
	"fmt"
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
			return fmt.Errorf("rollback for '%w' failed with: %v", err, ex)
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
		if txOpts == nil {
			txOpts = opt
		} else {
			txOpts.ReadOnly = opt.ReadOnly
			txOpts.Isolation = opt.Isolation
		}
	}
	tx, err := db.BeginTx(ctx, txOpts)
	if err != nil {
		return err
	}
	err = callback(tx)
	if err != nil {
		if ex := tx.Rollback(); ex != nil {
			return fmt.Errorf("rollback for '%w' failed with: %v", err, ex)
		}
		return err
	}
	return tx.Commit()
}
