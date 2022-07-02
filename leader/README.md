# leader

Good-enough [leader election](https://aws.amazon.com/builders-library/leader-election-in-distributed-systems/)
pattern implementation using MySQL as the coordinator.

Builds on the ideas found in [this gist](https://gist.github.com/ljjjustin/f2213ac9b9b8c31df746f8b56095ea32).

```go
import (
	"context"
	"database/sql"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/tomcz/gotools/leader"
)

func main() {
	var db *sql.DB // initialisation omitted

	node := leader.NewMysqlLeader(db, "app_leader")

	ctx, cancelElections := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		err := node.Acquire(ctx)
		if errors.Is(err, context.Canceled) {
			log.Println("elections canceled")
		} else {
			log.Println("elections failed:", err)
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < 100; i++ {
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
		cancelElections()
	}()

	wg.Wait()
}
```
