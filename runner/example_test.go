package runner_test

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/tomcz/gotools/runner"
)

var server *http.Server

func Example() {
	// create runner with default settings
	app := runner.New()

	// register a graceful shutdown
	app.ShutdownTimeout(server.Shutdown, 100*time.Millisecond)

	// kick off the server
	app.Run(func() error {
		err := server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	})

	// wait for SIGINT or SIGTERM signal from the os
	if err := app.Wait(); err != nil {
		log.Fatalln(err)
	}
}
