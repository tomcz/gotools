package leader

import "context"

type Leader interface {
	IsLeader(ctx context.Context) (bool, error)
	StartElections(ctx context.Context) error
}
