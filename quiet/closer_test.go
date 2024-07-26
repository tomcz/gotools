package quiet

import (
	"context"
	"testing"
	"time"

	assert "github.com/stretchr/testify/require"
)

func TestCloserOrder(t *testing.T) {
	var called []string
	c1 := func() {
		called = append(called, "1")
	}
	c2 := func() {
		called = append(called, "2")
	}
	c3 := func() {
		called = append(called, "3")
	}
	c4 := func() error {
		called = append(called, "4")
		return nil
	}
	c5 := func(context.Context) error {
		called = append(called, "5")
		return nil
	}

	c1w := &quietCloser{close: c1}
	c2w := &quietCloser{close: c2}

	closer := &Closer{}
	closer.Add(c1w, c2w)
	closer.AddFunc(c3)
	closer.AddFuncE(c4)
	closer.AddShutdown(c5, time.Minute)

	assert.NoError(t, closer.Close())
	assert.Equal(t, []string{"5", "4", "3", "2", "1"}, called)
}
