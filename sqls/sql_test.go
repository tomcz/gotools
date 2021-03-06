//go:build integration

package sqls

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const createLeaderTableSQL = `
CREATE TABLE IF NOT EXISTS test_leaders (
  id          int unsigned NOT NULL AUTO_INCREMENT,
  leader_name varchar(255) NOT NULL,
  node_name   varchar(255) NOT NULL,
  last_update datetime     NOT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY  unique_leader_name (leader_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
`

const insertLeaderSQL = `
INSERT INTO test_leaders (leader_name, node_name, last_update)
VALUES (?, ?, ?)
`

const selectLeadersSQL = `
SELECT leader_name, node_name
FROM test_leaders
`

const countLeadersSQL = `
SELECT count(*)
FROM test_leaders
WHERE leader_name = ?
`

func TestSqlTools(t *testing.T) {
	cfg := mysql.NewConfig()
	cfg.Net = "tcp"
	cfg.Addr = os.Getenv("DB_HOST")
	cfg.DBName = os.Getenv("DB_DATABASE")
	cfg.User = os.Getenv("DB_USER")
	cfg.Passwd = os.Getenv("DB_PASSWORD")

	db, err := sql.Open("mysql", cfg.FormatDSN())
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec(createLeaderTableSQL)
	require.NoError(t, err)

	tests := []struct {
		name   string
		testFn func(t *testing.T, db *sql.DB)
	}{
		{
			name:   "testInTxCommit",
			testFn: testInTxCommit,
		},
		{
			name:   "testInTxRollback",
			testFn: testInTxRollback,
		},
		{
			name:   "testInTxContextCommit",
			testFn: testInTxContextCommit,
		},
		{
			name:   "testInTxContextRollback",
			testFn: testInTxContextRollback,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFn(t, db)
		})
	}
}

func testInTxCommit(t *testing.T, db *sql.DB) {
	leaders := []struct {
		leader string
		node   string
	}{
		{
			leader: uuid.New().String(),
			node:   uuid.New().String(),
		},
		{
			leader: uuid.New().String(),
			node:   uuid.New().String(),
		},
	}
	err := InTx(db, func(tx *sql.Tx) error {
		for _, x := range leaders {
			_, err := tx.Exec(insertLeaderSQL, x.leader, x.node, time.Now())
			if err != nil {
				return err
			}
		}
		return nil
	})
	if assert.NoError(t, err) {
		results := make(map[string]string)
		err = QueryRows(db, selectLeadersSQL)(func(row ScanFunc) error {
			var leader, node string
			if err := row(&leader, &node); err != nil {
				return nil
			}
			results[leader] = node
			return nil
		})
		if assert.NoError(t, err) {
			assert.NotEmpty(t, results)
			for _, x := range leaders {
				assert.Equal(t, x.node, results[x.leader])
			}
		}
	}
}

func testInTxRollback(t *testing.T, db *sql.DB) {
	leaderName := uuid.New().String()
	leaders := []struct {
		leader string
		node   string
	}{
		{
			leader: leaderName,
			node:   uuid.New().String(),
		},
		{
			leader: leaderName,
			node:   uuid.New().String(),
		},
	}
	err := InTx(db, func(tx *sql.Tx) error {
		for _, x := range leaders {
			_, err := tx.Exec(insertLeaderSQL, x.leader, x.node, time.Now())
			if err != nil {
				return err
			}
		}
		return nil
	})
	if assert.Error(t, err) {
		var count int
		err = db.QueryRow(countLeadersSQL, leaderName).Scan(&count)
		if assert.NoError(t, err) {
			assert.Equal(t, 0, count)
		}
	}
}

func testInTxContextCommit(t *testing.T, db *sql.DB) {
	ctx := context.Background()
	leaders := []struct {
		leader string
		node   string
	}{
		{
			leader: uuid.New().String(),
			node:   uuid.New().String(),
		},
		{
			leader: uuid.New().String(),
			node:   uuid.New().String(),
		},
	}
	err := InTxContext(ctx, db, func(tx *sql.Tx) error {
		for _, x := range leaders {
			_, err := tx.ExecContext(ctx, insertLeaderSQL, x.leader, x.node, time.Now())
			if err != nil {
				return err
			}
		}
		return nil
	})
	if assert.NoError(t, err) {
		results := make(map[string]string)
		err = QueryRowsContext(ctx, db, selectLeadersSQL)(func(row ScanFunc) error {
			var leader, node string
			if err := row(&leader, &node); err != nil {
				return nil
			}
			results[leader] = node
			return nil
		})
		if assert.NoError(t, err) {
			assert.NotEmpty(t, results)
			for _, x := range leaders {
				assert.Equal(t, x.node, results[x.leader])
			}
		}
	}
}

func testInTxContextRollback(t *testing.T, db *sql.DB) {
	ctx := context.Background()
	leaderName := uuid.New().String()
	leaders := []struct {
		leader string
		node   string
	}{
		{
			leader: leaderName,
			node:   uuid.New().String(),
		},
		{
			leader: leaderName,
			node:   uuid.New().String(),
		},
	}
	err := InTxContext(ctx, db, func(tx *sql.Tx) error {
		for _, x := range leaders {
			_, err := tx.ExecContext(ctx, insertLeaderSQL, x.leader, x.node, time.Now())
			if err != nil {
				return err
			}
		}
		return nil
	})
	if assert.Error(t, err) {
		var count int
		err = db.QueryRowContext(ctx, countLeadersSQL, leaderName).Scan(&count)
		if assert.NoError(t, err) {
			assert.Equal(t, 0, count)
		}
	}
}
