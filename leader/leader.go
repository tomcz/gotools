package leader

import "context"

// Leader election, in distributed computing, is the process of designating
// a single process as the organizer of some task distributed among several
// computers (nodes). Before the task has begun, all network nodes are either
// unaware which node will serve as the "leader" (or coordinator) of the task,
// or unable to communicate with the current coordinator. After a leader election
// algorithm has been run, however, each node throughout the network recognizes
// a particular, unique node as the task leader.
type Leader interface {
	// IsLeader returns whether this node is the leader, or an error if it was
	// unable to determine if it is the leader for any reason.
	IsLeader(ctx context.Context) (bool, error)
	// RunElections is a blocking call. It should exit when the context is
	// cancelled, or when an error occurs during the election and the OnError
	// strategy indicates an end to elections.
	RunElections(ctx context.Context, onError OnError) error
}
