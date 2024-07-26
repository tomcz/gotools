package errgroup

import (
	"context"
	"fmt"
	"testing"

	assert "github.com/stretchr/testify/require"
)

func TestGroup_OK(t *testing.T) {
	group := New()
	group.Go(func() error {
		return nil
	})
	assert.NoError(t, group.Wait())
}

func TestGroup_Error(t *testing.T) {
	group, ctx := NewContext(context.Background())
	group.Go(func() error {
		return fmt.Errorf("oops")
	})
	group.Go(func() error {
		<-ctx.Done()
		return ctx.Err()
	})
	err := group.Wait()
	assert.Error(t, err, "oops")
}

func TestGroup_Panic(t *testing.T) {
	group, ctx := NewContext(context.Background())
	group.Go(func() error {
		panic("doh")
	})
	group.Go(func() error {
		<-ctx.Done()
		return ctx.Err()
	})
	err := group.Wait()
	assert.ErrorContains(t, err, "panic: doh")
}

func TestGroup_Panic_Error(t *testing.T) {
	cause := fmt.Errorf("frack")
	group, ctx := NewContext(context.Background())
	group.Go(func() error {
		panic(cause)
	})
	group.Go(func() error {
		<-ctx.Done()
		return ctx.Err()
	})
	err := group.Wait()
	assert.ErrorIs(t, err, cause)
}

func TestGroup_Panic_Handler(t *testing.T) {
	ph := func(p any) error {
		return fmt.Errorf("%v handled", p)
	}
	group := New(WithPanicHandler(ph))
	group.Go(func() error {
		panic("mischief")
	})
	assert.Error(t, group.Wait(), "mischief handled")
}
