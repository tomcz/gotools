//go:build integration

package leader

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"gotest.tools/v3/assert"
)

func TestMysqlLeader(t *testing.T) {
	cfg := mysql.NewConfig()
	cfg.Net = "tcp"
	cfg.Addr = os.Getenv("DB_HOST")
	cfg.DBName = os.Getenv("DB_DATABASE")
	cfg.User = os.Getenv("DB_USER")
	cfg.Passwd = os.Getenv("DB_PASSWORD")

	db, err := sql.Open("mysql", cfg.FormatDSN())
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
	leaderName := uuid.New().String()
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
		leaderName: uuid.New().String(),
		nodeName:   uuid.New().String(),
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
	leaderName := uuid.New().String()

	mock := clock.NewMock()
	mock.Set(time.Now())

	l1 := &mysqlLeader{
		db:         db,
		leaderName: leaderName,
		nodeName:   uuid.New().String(),
		clock:      mock,
		age:        10 * time.Second,
	}
	l2 := &mysqlLeader{
		db:         db,
		leaderName: leaderName,
		nodeName:   uuid.New().String(),
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
