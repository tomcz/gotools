package runner

import (
	"context"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tomcz/gotools/errgroup"
	"github.com/tomcz/gotools/quiet"
)

type cfgOpt struct {
	handler errgroup.PanicHandler
	logger  quiet.Logger
	signals []os.Signal
	async   bool
}

// Opt is a Runner configuration option.
type Opt func(cfg *cfgOpt)

// WithHandler overrides the default panic handler.
func WithHandler(handler errgroup.PanicHandler) Opt {
	return func(cfg *cfgOpt) {
		cfg.handler = handler
	}
}

// WithLogger overrides the default shutdown logger.
func WithLogger(logger quiet.Logger) Opt {
	return func(cfg *cfgOpt) {
		cfg.logger = logger
	}
}

// WithSignals overrides the default [syscall.SIGINT] and [syscall.SIGTERM] shutdown signals.
func WithSignals(signals ...os.Signal) Opt {
	return func(cfg *cfgOpt) {
		cfg.signals = signals
	}
}

// WithAsyncCleanup will run the registered cleanup functions in parallel rather than the default
// inverse order of registration that mimics the invocation order of defers in a function.
func WithAsyncCleanup() Opt {
	return func(cfg *cfgOpt) {
		cfg.async = true
	}
}

// Runner will kick off all runner functions passed to the [Runner.Run] function.
// It will invoke all registered cleanup functions when any of the runners terminates
// or when an exit signal (i.e. [syscall.SIGINT] or [syscall.SIGTERM]) is received.
type Runner struct {
	ctx     context.Context
	cancel  context.CancelFunc
	group   *errgroup.PanicGroup
	closer  *quiet.Closer
	signals []os.Signal
	async   bool
}

// New constructor.
func New(opts ...Opt) *Runner {
	cfg := &cfgOpt{}
	for _, opt := range opts {
		opt(cfg)
	}

	ctx, cancel := context.WithCancel(context.Background())

	group := errgroup.New()
	if cfg.handler != nil {
		group.SetPanicHandler(cfg.handler)
	}

	closer := &quiet.Closer{}
	if cfg.logger != nil {
		closer.SetLogger(cfg.logger)
	}

	signals := []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	if len(cfg.signals) > 0 {
		signals = cfg.signals
	}

	return &Runner{
		ctx:     ctx,
		cancel:  cancel,
		group:   group,
		closer:  closer,
		signals: signals,
		async:   cfg.async,
	}
}

// Cleanup registration.
func (r *Runner) Cleanup(closers ...io.Closer) {
	r.closer.Add(closers...)
}

// CleanupFunc registration.
func (r *Runner) CleanupFunc(closer func()) {
	r.closer.AddFunc(closer)
}

// CleanupFuncE registration.
func (r *Runner) CleanupFuncE(closer func() error) {
	r.closer.AddFuncE(closer)
}

// CleanupTimeout registration.
func (r *Runner) CleanupTimeout(closer func(ctx context.Context) error, timeout time.Duration) {
	r.closer.AddTimeout(closer, timeout)
}

// Run the given runner function in the background.
func (r *Runner) Run(fn func() error) {
	r.group.Go(func() error {
		defer r.cancel()
		return fn()
	})
}

// Wait for the runners to terminate. Kick off the registered cleanup functions when any runner
// terminates or when an exit signal (i.e. [syscall.SIGINT] or [syscall.SIGTERM]) is received.
func (r *Runner) Wait() error {
	r.group.Go(func() error {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, r.signals...)
		select {
		case <-sigint:
			return r.cleanup()
		case <-r.ctx.Done():
			return r.cleanup()
		}
	})
	return r.group.Wait()
}

func (r *Runner) cleanup() error {
	if r.async {
		r.closer.CloseAsync()
	} else {
		r.closer.CloseAll()
	}
	return nil
}
