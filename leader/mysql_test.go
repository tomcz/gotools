//go:build integration

package leader

import (
	"context"
	"database/sql"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"gotest.tools/v3/assert"
)

func openConn(dbName string) (*sql.DB, error) {
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
	db, err := openConn("")
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
	db, err := openConn("")
	if err != nil {
		return err
	}
	_, err = db.Exec("DROP DATABASE " + dbName)
	if err != nil {
		return err
	}
	return db.Close()
}

func TestMysqlLeader(t *testing.T) {
	dbName, err := createTestDatabase()
	assert.NilError(t, err, "createTestDatabase failed")
	defer func() {
		dropErr := dropTestDatabase(dbName)
		assert.Check(t, dropErr, "dropTestDatabase failed")
	}()

	db, err := openConn(dbName)
	assert.NilError(t, err)
	defer db.Close()

	assert.NilError(t, CreateMysqlLeaderTable(db))

	tests := []struct {
		name   string
		testFn func(t *testing.T, db *sql.DB)
	}{
		{
			name:   "testUnelectedIsNotLeader",
			testFn: testUnelectedIsNotLeader,
		},
		{
			name:   "testElectedIsLeader",
			testFn: testElectedIsLeader,
		},
		{
			name:   "testElectionWinnerIsLeader",
			testFn: testElectionWinnerIsLeader,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFn(t, db)
		})
	}
}

func testUnelectedIsNotLeader(t *testing.T, db *sql.DB) {
	ctx := context.Background()
	leaderName := uuid.NewString()
	leader := NewMysqlLeader(db, leaderName)
	isLeader, err := leader.IsLeader(ctx)
	assert.NilError(t, err)
	assert.Assert(t, !isLeader)
}

func testElectedIsLeader(t *testing.T, db *sql.DB) {
	ctx := context.Background()

	mock := clock.NewMock()
	mock.Set(time.Now())

	leader := &mysqlLeader{
		db:         db,
		leaderName: uuid.NewString(),
		nodeName:   uuid.NewString(),
		clock:      mock,
		age:        10 * time.Second,
	}

	assert.NilError(t, leader.election(ctx))

	isLeader, err := leader.IsLeader(ctx)
	assert.NilError(t, err)
	assert.Assert(t, isLeader)
}

func testElectionWinnerIsLeader(t *testing.T, db *sql.DB) {
	ctx := context.Background()
	leaderName := uuid.NewString()

	mock := clock.NewMock()
	mock.Set(time.Now())

	l1 := &mysqlLeader{
		db:         db,
		leaderName: leaderName,
		nodeName:   uuid.NewString(),
		clock:      mock,
		age:        10 * time.Second,
	}
	l2 := &mysqlLeader{
		db:         db,
		leaderName: leaderName,
		nodeName:   uuid.NewString(),
		clock:      mock,
		age:        10 * time.Second,
	}

	// l1 wins election
	assert.NilError(t, l1.election(ctx))
	assert.NilError(t, l2.election(ctx))

	isLeader, err := l1.IsLeader(ctx)
	assert.NilError(t, err)
	assert.Assert(t, isLeader, "l1 should be leader")

	isLeader, err = l2.IsLeader(ctx)
	assert.NilError(t, err)
	assert.Assert(t, !isLeader, "l2 should not be leader")

	// l2 wins next election
	mock.Add(11 * time.Second)
	assert.NilError(t, l2.election(ctx))
	assert.NilError(t, l1.election(ctx))

	isLeader, err = l1.IsLeader(ctx)
	assert.NilError(t, err)
	assert.Assert(t, !isLeader, "l1 should not be leader")

	isLeader, err = l2.IsLeader(ctx)
	assert.NilError(t, err)
	assert.Assert(t, isLeader, "l2 should be leader")
}
