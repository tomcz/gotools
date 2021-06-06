package transaction

import (
	"context"
	"database/sql"
	"fmt"
)

type TxFn func(tx *sql.Tx) error

type TxCtxFn func(ctx context.Context, tx *sql.Tx) error

func InTx(db *sql.DB, fn TxFn) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	err = fn(tx)
	if err != nil {
		if ex := tx.Rollback(); ex != nil {
			return fmt.Errorf("rollback for '%w' failed with: %v", err, ex)
		}
		return err
	}
	return tx.Commit()
}

func InTxContext(ctx context.Context, db *sql.DB, fn TxCtxFn, opts ...*sql.TxOptions) error {
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
	err = fn(ctx, tx)
	if err != nil {
		if ex := tx.Rollback(); ex != nil {
			return fmt.Errorf("rollback for '%w' failed with: %v", err, ex)
		}
		return err
	}
	return tx.Commit()
}
