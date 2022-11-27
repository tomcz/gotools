package errgroup

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroup_OK(t *testing.T) {
	group := New()
	group.Go(func() error {
		return nil
	})
	assert.NoError(t, group.Wait())
}

func TestGroup_Error(t *testing.T) {
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

func TestGroup_Panic(t *testing.T) {
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

func TestGroup_Panic_Error(t *testing.T) {
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

func TestGroup_Panic_Handler(t *testing.T) {
	ph := func(p any) error {
		return fmt.Errorf("%v handled", p)
	}
	group := New()
	group.OnPanic(ph)
	group.Go(func() error {
		panic("mischief")
	})
	assert.EqualError(t, group.Wait(), "mischief handled")
}
