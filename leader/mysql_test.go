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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMysqlLeader(t *testing.T) {
	cfg := mysql.NewConfig()
	cfg.Addr = os.Getenv("DB_HOST")
	cfg.DBName = os.Getenv("DB_DATABASE")
	cfg.User = os.Getenv("DB_USER")
	cfg.Passwd = os.Getenv("DB_PASSWORD")

	db, err := sql.Open("mysql", cfg.FormatDSN())
	require.NoError(t, err)
	defer db.Close()

	require.NoError(t, CreateMysqlLeaderTable(db))

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
	if assert.NoError(t, err) {
		assert.False(t, isLeader)
	}
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

	election := leader.election()
	require.NoError(t, election(ctx))

	isLeader, err := leader.IsLeader(ctx)
	if assert.NoError(t, err) {
		assert.True(t, isLeader)
	}
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

	e1 := l1.election()
	e2 := l2.election()

	// l1 wins election
	require.NoError(t, e1(ctx))
	require.NoError(t, e2(ctx))

	isLeader, err := l1.IsLeader(ctx)
	if assert.NoError(t, err) {
		assert.True(t, isLeader, "l1 should be leader")
	}
	isLeader, err = l2.IsLeader(ctx)
	if assert.NoError(t, err) {
		assert.False(t, isLeader, "l2 should not be leader")
	}

	// l2 wins next election
	mock.Add(11 * time.Second)
	require.NoError(t, e2(ctx))
	require.NoError(t, e1(ctx))

	isLeader, err = l1.IsLeader(ctx)
	if assert.NoError(t, err) {
		assert.False(t, isLeader, "l1 should not be leader")
	}
	isLeader, err = l2.IsLeader(ctx)
	if assert.NoError(t, err) {
		assert.True(t, isLeader, "l2 should be leader")
	}
}
