package errgroup

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroupOK(t *testing.T) {
	group := New()
	group.Go(func() error {
		return nil
	})
	assert.NoError(t, group.Wait())
}

func TestGroupError(t *testing.T) {
	group, ctx := WithContext(context.Background())
	group.Go(func() error {
		return fmt.Errorf("oops")
	})
	group.Go(func() error {
		<-ctx.Done()
		return ctx.Err()
	})
	err := group.Wait()
	assert.EqualError(t, err, "oops")
}

func TestGroupPanic(t *testing.T) {
	group, ctx := WithContext(context.Background())
	group.Go(func() error {
		panic("doh")
	})
	group.Go(func() error {
		<-ctx.Done()
		return ctx.Err()
	})
	err := group.Wait()
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "panic: doh")
	}
}

func TestGroupPanicError(t *testing.T) {
	cause := fmt.Errorf("frack")
	group, ctx := WithContext(context.Background())
	group.Go(func() error {
		panic(cause)
	})
	group.Go(func() error {
		<-ctx.Done()
		return ctx.Err()
	})
	err := group.Wait()
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, cause)
	}
}
