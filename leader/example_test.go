package leader_test

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/tomcz/gotools/leader"
)

var (
	cc context.Context
	db *sql.DB
)

func Example() {
	ctx, cancel := context.WithCancel(cc)

	node := leader.NewMysqlLeader(db, "app_leader")

	var wg sync.WaitGroup
	wg.Go(func() {
		// node.Acquire is a blocking call
		err := node.Acquire(ctx)
		if errors.Is(err, context.Canceled) {
			log.Println("elections canceled")
		} else {
			log.Println("elections failed:", err)
		}
	})
	wg.Go(func() {
		for i := range 100 {
			if i > 0 {
				time.Sleep(time.Second)
			}
			isLeader, err := node.IsLeader(ctx)
			if err != nil {
				log.Println("leader check failed:", err)
				continue
			}
			if isLeader {
				log.Println("I am the Leader :)")
				continue
			}
			log.Println("I am NOT the Leader :(")
		}
		// unblock node.Acquire
		cancel()
	})
	wg.Wait()
}
