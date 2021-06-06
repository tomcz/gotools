package leader

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/google/uuid"
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

type mysqlLeader struct {
	db         *sql.DB
	leaderName string
	nodeName   string
	clock      clock.Clock
	tick       time.Duration
	age        time.Duration
}

// NewMysqlLeader provides an implementation of the Leader interface using
// MySQL as the point of coordination between nodes. It is not a perfect
// leadership election implementation but should be good enough providing
// that tasks that require leadership election do not run for longer than
// either the tick or age intervals.
func NewMysqlLeader(db *sql.DB, leaderName string, opts ...MysqlOpt) Leader {
	leader := &mysqlLeader{
		db:         db,
		leaderName: leaderName,
		clock:      clock.New(),
	}
	for _, opt := range opts {
		opt(leader)
	}
	if leader.nodeName == "" {
		leader.nodeName = uuid.New().String()
	}
	if leader.tick < time.Second {
		leader.tick = 15 * time.Second
	}
	if leader.age < time.Second {
		leader.age = 60 * time.Second
	}
	return leader
}

func (m *mysqlLeader) StartElections(ctx context.Context) error {
	election := m.election()
	ticker := m.clock.Ticker(m.tick)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := election(ctx); err != nil {
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
node_name = IF(last_update < DATE_SUB(VALUES(last_update), INTERVAL %d SECOND), VALUES(node_name), node_name),
last_update = IF(node_name = VALUES(node_name), VALUES(last_update), last_update)
`

func (m *mysqlLeader) election() func(context.Context) error {
	stmt := fmt.Sprintf(electionSQL, int64(m.age.Seconds()))
	return func(ctx context.Context) error {
		_, err := m.db.ExecContext(ctx, stmt, m.leaderName, m.nodeName, m.clock.Now())
		return err
	}
}

// CreateMysqlLeaderSQL is the create statement used by CreateMysqlLeaderTable.
// It's published so that it can be used in database migrations without needing
// to call the CreateMysqlLeaderTable function.
const CreateMysqlLeaderSQL = `
CREATE TABLE IF NOT EXISTS leader_election (
  id          int unsigned NOT NULL AUTO_INCREMENT,
  leader_name varchar(255) NOT NULL,
  node_name   varchar(255) NOT NULL,
  last_update datetime     NOT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY  unique_leader_name (leader_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
`

// CreateMysqlLeaderTable sets up the leadership election table and its constraints.
// It is not part of the MysqlLeader object since in practice it's a bad idea to run
// services with permissions to create or modify database schemas.
func CreateMysqlLeaderTable(db *sql.DB) error {
	_, err := db.Exec(CreateMysqlLeaderSQL)
	return err
}
