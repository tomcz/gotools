package leader

import "context"

// Leader election, In distributed computing, is the process of designating
// a single process as the organizer of some task distributed among several
// computers (nodes). Before the task has began, all network nodes are either
// unaware which node will serve as the "leader" (or coordinator) of the task,
// or unable to communicate with the current coordinator. After a leader election
// algorithm has been run, however, each node throughout the network recognizes
//  a particular, unique node as the task leader.
type Leader interface {
	// IsLeader returns whether this node is the leader, or an error if it was
	// unable to determine if it is the leader for any reason.
	IsLeader(ctx context.Context) (bool, error)
	// StartElections is a blocking call. It should exit when the context is
	// cancelled or if there is an error in running the election algorithm.
	StartElections(ctx context.Context) error
}
