package leader

import "context"

type alwaysLeader struct{}

// NewAlwaysLeader returns an implementation of the Leader interface
// that always considers itself to be the leader.
func NewAlwaysLeader() Leader {
	return alwaysLeader{}
}

func (a alwaysLeader) IsLeader(ctx context.Context) (bool, error) {
	return true, nil
}

func (a alwaysLeader) StartElections(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}
