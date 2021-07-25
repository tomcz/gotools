package leader

import "log"

// OnError provides the error handling strategy for Leader.RunElections.
// When an election error is encountered it will be passed to OnError.
// If OnError returns a non-nil error the Leader.RunElections blocking
// call will exit with the error returned by OnError, signifying an end
// to elections.
type OnError func(err error) error

// StopOnError will terminate elections on any error.
func StopOnError(err error) error {
	return err
}

// ContinueOnError allows elections to continue on any error.
// The error iself is logged by the standard library logger.
func ContinueOnError(err error) error {
	log.Println("election error:", err)
	return nil
}
