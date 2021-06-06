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
