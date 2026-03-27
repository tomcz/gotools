# runner

Combines an errgoup.Group, a quiet.Closer, and an os/signal receiver to run functions in the background, wait for a termination signal, and then run all given cleanup tasks.

```go
import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/tomcz/gotools/runner"
)

func main() {
	server := &http.Server{} // TODO

	// create runner with default settings
	app := runner.New()

	// kick off the server
	app.Run(func() error {
		err := server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	})

	// register a graceful shutdown
	app.CleanupTimeout(server.Shutdown, 100*time.Millisecond)

	// wait for SIGINT or SIGTERM signal from the os
	if err := app.Wait(); err != nil {
		log.Fatalln(err)
	}
}
```
