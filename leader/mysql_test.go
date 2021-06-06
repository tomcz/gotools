package leader

import (
	"context"
	"database/sql"
	"errors"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mock := clock.NewMock()
	leaderName := uuid.New().String()
	leader := NewMysqlLeader(db, leaderName, WithClock(mock), WithTick(time.Second))

	go func() {
		if err := leader.StartElections(ctx); errors.Is(err, context.Canceled) {
			t.Log("unexpected election error:", err)
		}
	}()
	mock.Add(time.Second)

	isLeader, err := leader.IsLeader(ctx)
	if assert.NoError(t, err) {
		assert.True(t, isLeader)
	}
}
