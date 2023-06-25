//go:build integration

package txx

import (
	"context"
	"database/sql"
	"os"
	"strings"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"gotest.tools/v3/assert"
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
VALUES (:leader, :node, NOW())
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
		assert.Check(t, dropErr, "dropTestDatabase failed")
	}()

	db, err := sqlOpen(dbName)
	assert.NilError(t, err)
	defer db.Close()

	_, err = db.Exec(createLeaderTableSQL)
	assert.NilError(t, err)

	tests := []struct {
		name   string
		testFn func(t *testing.T, db *sqlx.DB)
	}{
		{
			name:   "testInTxxCommit",
			testFn: testInTxxCommit,
		},
		{
			name:   "testInTxxRollback",
			testFn: testInTxxRollback,
		},
	}
	dbx := sqlx.NewDb(db, "mysql")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFn(t, dbx)
		})
	}
}

type testLeader struct {
	Name string `db:"leader_name"`
	Node string `db:"node_name"`
}

func testInTxxCommit(t *testing.T, db *sqlx.DB) {
	ctx := context.Background()
	leaders := map[string]string{
		uuid.NewString(): uuid.NewString(),
		uuid.NewString(): uuid.NewString(),
	}
	err := InTxx(ctx, db, func(tx *sqlx.Tx) error {
		for leader, node := range leaders {
			data := map[string]any{"leader": leader, "node": node}
			_, rerr := tx.NamedExecContext(ctx, insertLeaderSQL, data)
			if rerr != nil {
				return rerr
			}
		}
		return nil
	})
	assert.NilError(t, err)

	var results []testLeader
	err = db.SelectContext(ctx, &results, selectLeadersSQL)
	assert.NilError(t, err)
	for _, result := range results {
		assert.Equal(t, leaders[result.Name], result.Node)
	}
}

func testInTxxRollback(t *testing.T, db *sqlx.DB) {
	ctx := context.Background()
	leaderName := uuid.NewString()
	leaders := []map[string]any{
		{"leader": leaderName, "node": uuid.NewString()},
		{"leader": leaderName, "node": uuid.NewString()},
	}
	err := InTxx(ctx, db, func(tx *sqlx.Tx) error {
		for _, data := range leaders {
			_, rerr := tx.NamedExecContext(ctx, insertLeaderSQL, data)
			if rerr != nil {
				return rerr
			}
		}
		return nil
	})
	assert.ErrorContains(t, err, "Duplicate entry")

	var count int
	err = db.GetContext(ctx, &count, countLeadersSQL, leaderName)
	assert.NilError(t, err)
	assert.Equal(t, 0, count)
}
