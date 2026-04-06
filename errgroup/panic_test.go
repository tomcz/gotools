package errgroup

import (
	"context"
	"fmt"
	"testing"

	"gotest.tools/v3/assert"
)

func TestGroup_OK(t *testing.T) {
	group := New()
	group.Go(func() error {
		return nil
	})
	assert.NilError(t, group.Wait())
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
		panic("d'oh")
	})
	group.Go(func() error {
		<-ctx.Done()
		return ctx.Err()
	})
	err := group.Wait()
	assert.ErrorContains(t, err, "panic: d'oh")
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

type testPanicHandler struct{}

func (h testPanicHandler) Panic(p any) error {
	return fmt.Errorf("%v crap", p)
}

func TestGroup_Panic_Handler(t *testing.T) {
	group := New()
	group.SetPanicHandler(testPanicHandler{})
	group.Go(func() error {
		panic("d'oh")
	})
	assert.Error(t, group.Wait(), "d'oh crap")
}
