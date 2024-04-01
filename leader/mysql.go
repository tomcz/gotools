package leader

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	"k8s.io/utils/clock"
)

// MysqlOpt allows configuration of leader defaults.
type MysqlOpt func(leader *mysqlLeader)

// WithNodeName allows the node name to be specified.
// The default value is a random UUID.
func WithNodeName(name string) MysqlOpt {
	return func(leader *mysqlLeader) {
		leader.nodeName = name
	}
}

// WithTick allows the default election frequency to
// be specified. The default is 15 seconds.
func WithTick(tick time.Duration) MysqlOpt {
	return func(leader *mysqlLeader) {
		leader.tick = tick
	}
}

// WithAge allows the default lifespan of an election
// to be specified. The default is 60 seconds.
func WithAge(age time.Duration) MysqlOpt {
	return func(leader *mysqlLeader) {
		leader.age = age
	}
}

// WithClock replaces the default system clock.
func WithClock(ck clock.WithTicker) MysqlOpt {
	return func(leader *mysqlLeader) {
		leader.clock = ck
	}
}

// WithOnError allows the default strategy of
// terminating leadership elections on errors
// during the Leader.Acquire blocking call to
// be replaced with something more nuanced.
// If the onError strategy returns a non-nil
// error value, the blocking call will exit
// with the returned error. If the strategy
// returns a nil value, leadership election
// will continue on the next clock tick.
func WithOnError(onError func(error) error) MysqlOpt {
	return func(leader *mysqlLeader) {
		leader.onError = onError
	}
}

// AbortOnError is the default WithOnError strategy.
func AbortOnError(err error) error {
	return err
}

// ContinueOnError is an example WithOnError strategy that logs
// the error and allows the leadership election to proceed.
func ContinueOnError(err error) error {
	log.Printf("[WARNING] leadership error: %v\n", err)
	return nil
}

type mysqlLeader struct {
	db         *sql.DB
	leaderName string
	nodeName   string
	clock      clock.WithTicker
	tick       time.Duration
	age        time.Duration
	onError    func(error) error
}

// NewMysqlLeader provides an implementation of the Leader interface using
// MySQL as the point of coordination between nodes. It is not a perfect
// leadership election implementation but should be good enough providing
// that tasks that require leadership election do not run for longer than
// either the tick or age intervals.
func NewMysqlLeader(db *sql.DB, leaderName string, opts ...MysqlOpt) Leader {
	leader := &mysqlLeader{db: db, leaderName: leaderName}
	for _, opt := range opts {
		opt(leader)
	}
	if leader.nodeName == "" {
		leader.nodeName = uuid.NewString()
	}
	if leader.clock == nil {
		leader.clock = clock.RealClock{}
	}
	if leader.tick < time.Second {
		leader.tick = 15 * time.Second
	}
	if leader.age < time.Second {
		leader.age = 60 * time.Second
	}
	if leader.onError == nil {
		leader.onError = AbortOnError
	}
	return leader
}

func (m *mysqlLeader) Acquire(ctx context.Context) error {
	ticker := m.clock.NewTicker(m.tick)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C():
			if err := m.election(ctx); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

const isLeaderSQL = `
SELECT COUNT(*)
FROM leader_election
WHERE leader_name = ?
AND node_name = ?
`

func (m *mysqlLeader) IsLeader(ctx context.Context) (bool, error) {
	var count int
	err := m.db.QueryRowContext(ctx, isLeaderSQL, m.leaderName, m.nodeName).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

const electionSQL = `
INSERT INTO leader_election (leader_name, node_name, last_update) VALUES (?, ?, ?)
ON DUPLICATE KEY UPDATE
node_name = IF(last_update < DATE_SUB(VALUES(last_update), INTERVAL ? SECOND), VALUES(node_name), node_name),
last_update = IF(node_name = VALUES(node_name) OR last_update < DATE_SUB(VALUES(last_update), INTERVAL ? SECOND), VALUES(last_update), last_update)
`

func (m *mysqlLeader) election(ctx context.Context) error {
	_, err := m.db.ExecContext(ctx, electionSQL, m.leaderName, m.nodeName, m.clock.Now(), int64(m.age.Seconds()), int64(m.age.Seconds()))
	if err != nil {
		err = m.onError(err)
	}
	return err
}

// CreateMysqlLeaderSQL is the create statement used by CreateMysqlLeaderTable.
// It's published so that it can be used in database migrations without needing
// to call the CreateMysqlLeaderTable function.
const CreateMysqlLeaderSQL = `
CREATE TABLE IF NOT EXISTS leader_election (
  leader_name varchar(255) NOT NULL PRIMARY KEY,
  node_name   varchar(255) NOT NULL,
  last_update datetime     NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
`

// CreateMysqlLeaderTable sets up the leadership election table and its constraints.
// It is not part of the MysqlLeader object since in practice it's a bad idea to run
// services with permissions to create or modify database schemas.
func CreateMysqlLeaderTable(db *sql.DB) error {
	_, err := db.Exec(CreateMysqlLeaderSQL)
	return err
}
