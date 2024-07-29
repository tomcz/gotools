//go:build integration

package sqls

import (
	"context"
	"database/sql"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
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

func sqlOpen(dbName string) (*sql.DB, error) {
	cfg := mysql.NewConfig()
	cfg.Net = "tcp"
	cfg.Addr = os.Getenv("DB_HOST")
	cfg.User = os.Getenv("DB_USER")
	cfg.Passwd = os.Getenv("DB_PASSWORD")
	cfg.DBName = dbName
	return sql.Open("mysql", cfg.FormatDSN())
}

func createTestDatabase() (string, error) {
	dbName := strings.ReplaceAll("test"+uuid.NewString(), "-", "")
	db, err := sqlOpen("")
	if err != nil {
		return "", err
	}
	_, err = db.Exec("CREATE DATABASE " + dbName)
	if err != nil {
		return "", err
	}
	return dbName, db.Close()
}

func dropTestDatabase(dbName string) error {
	db, err := sqlOpen("")
	if err != nil {
		return err
	}
	_, err = db.Exec("DROP DATABASE " + dbName)
	if err != nil {
		return err
	}
	return db.Close()
}

func TestSqlTools(t *testing.T) {
	dbName, err := createTestDatabase()
	assert.NilError(t, err, "createTestDatabase failed")
	defer func() {
		dropErr := dropTestDatabase(dbName)
		assert.Check(t, is.Nil(dropErr), "dropTestDatabase failed")
	}()

	db, err := sqlOpen(dbName)
	assert.NilError(t, err)
	defer db.Close()

	_, err = db.Exec(createLeaderTableSQL)
	assert.NilError(t, err)

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
			leader: uuid.NewString(),
			node:   uuid.NewString(),
		},
		{
			leader: uuid.NewString(),
			node:   uuid.NewString(),
		},
	}
	err := InTx(db, func(tx *sql.Tx) error {
		for _, x := range leaders {
			_, rerr := tx.Exec(insertLeaderSQL, x.leader, x.node, time.Now())
			if rerr != nil {
				return rerr
			}
		}
		return nil
	})
	assert.NilError(t, err)
	results := make(map[string]string)
	err = QueryRows(db, selectLeadersSQL)(func(row ScanFunc) error {
		var leader, node string
		if rerr := row(&leader, &node); rerr != nil {
			return rerr
		}
		results[leader] = node
		return nil
	})
	assert.NilError(t, err)
	assert.Assert(t, len(results) > 0)
	for _, x := range leaders {
		assert.Equal(t, x.node, results[x.leader])
	}
}

func testInTxRollback(t *testing.T, db *sql.DB) {
	leaderName := uuid.NewString()
	leaders := []struct {
		leader string
		node   string
	}{
		{
			leader: leaderName,
			node:   uuid.NewString(),
		},
		{
			leader: leaderName,
			node:   uuid.NewString(),
		},
	}
	err := InTx(db, func(tx *sql.Tx) error {
		for _, x := range leaders {
			_, rerr := tx.Exec(insertLeaderSQL, x.leader, x.node, time.Now())
			if rerr != nil {
				return rerr
			}
		}
		return nil
	})
	assert.ErrorContains(t, err, "Duplicate entry")
	var count int
	err = db.QueryRow(countLeadersSQL, leaderName).Scan(&count)
	assert.NilError(t, err)
	assert.Equal(t, 0, count)
}

func testInTxContextCommit(t *testing.T, db *sql.DB) {
	ctx := context.Background()
	leaders := []struct {
		leader string
		node   string
	}{
		{
			leader: uuid.NewString(),
			node:   uuid.NewString(),
		},
		{
			leader: uuid.NewString(),
			node:   uuid.NewString(),
		},
	}
	err := InTxContext(ctx, db, func(tx *sql.Tx) error {
		for _, x := range leaders {
			_, rerr := tx.ExecContext(ctx, insertLeaderSQL, x.leader, x.node, time.Now())
			if rerr != nil {
				return rerr
			}
		}
		return nil
	})
	assert.NilError(t, err)
	results := make(map[string]string)
	err = QueryRowsContext(ctx, db, selectLeadersSQL)(func(row ScanFunc) error {
		var leader, node string
		if rerr := row(&leader, &node); rerr != nil {
			return rerr
		}
		results[leader] = node
		return nil
	})
	assert.NilError(t, err)
	assert.Assert(t, len(results) > 0)
	for _, x := range leaders {
		assert.Equal(t, x.node, results[x.leader])
	}
}

func testInTxContextRollback(t *testing.T, db *sql.DB) {
	ctx := context.Background()
	leaderName := uuid.NewString()
	leaders := []struct {
		leader string
		node   string
	}{
		{
			leader: leaderName,
			node:   uuid.NewString(),
		},
		{
			leader: leaderName,
			node:   uuid.NewString(),
		},
	}
	err := InTxContext(ctx, db, func(tx *sql.Tx) error {
		for _, x := range leaders {
			_, rerr := tx.ExecContext(ctx, insertLeaderSQL, x.leader, x.node, time.Now())
			if rerr != nil {
				return rerr
			}
		}
		return nil
	})
	assert.ErrorContains(t, err, "Duplicate entry")
	var count int
	err = db.QueryRowContext(ctx, countLeadersSQL, leaderName).Scan(&count)
	assert.NilError(t, err)
	assert.Equal(t, 0, count)
}
