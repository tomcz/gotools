package leader

import "context"

type Leader interface {
	IsLeader(ctx context.Context) (bool, error)
	Election(ctx context.Context) error
}
