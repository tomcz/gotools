package sqls_test

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/tomcz/gotools/sqls"
)

var (
	ctx context.Context
	db  *sql.DB
)

func ExampleQueryRows() {
	selectSQL := "select leader, node from current_leaders"
	results := make(map[string]string)

	err := sqls.QueryRows(db, selectSQL)(func(row sqls.ScanFunc) error {
		var leader, node string
		if err := row(&leader, &node); err != nil {
			return nil
		}
		results[leader] = node
		return nil
	})

	if err != nil {
		log.Fatalln(err)
	}
	log.Println(results)
}

func ExampleQueryRowsContext() {
	selectSQL := "select leader, node from current_leaders"
	results := make(map[string]string)

	err := sqls.QueryRowsContext(ctx, db, selectSQL)(func(row sqls.ScanFunc) error {
		var leader, node string
		if err := row(&leader, &node); err != nil {
			return nil
		}
		results[leader] = node
		return nil
	})

	if err != nil {
		log.Fatalln(err)
	}
	log.Println(results)
}

func ExampleInTx() {
	inserSQL := "insert into old_leaders (leader, node, created_at) values (?, ?, ?)"
	leaders := map[string]string{"wibble": "foo", "wobble": "bar"}

	err := sqls.InTx(db, func(tx *sql.Tx) error {
		for leader, node := range leaders {
			_, err := tx.Exec(inserSQL, leader, node, time.Now())
			if err != nil {
				return err // tx will be rolled-back
			}
		}
		return nil // tx will be committed
	})

	if err != nil {
		log.Fatalln(err)
	}
}

func ExampleInTxContext() {
	inserSQL := "insert into old_leaders (leader, node, created_at) values (?, ?, ?)"
	leaders := map[string]string{"wibble": "foo", "wobble": "bar"}

	err := sqls.InTxContext(ctx, db, func(tx *sql.Tx) error {
		for leader, node := range leaders {
			_, err := tx.ExecContext(ctx, inserSQL, leader, node, time.Now())
			if err != nil {
				return err // tx will be rolled-back
			}
		}
		return nil // tx will be committed
	})

	if err != nil {
		log.Fatalln(err)
	}
}
