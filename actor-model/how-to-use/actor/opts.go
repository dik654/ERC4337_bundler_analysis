package actor

import (
	"context"
	"time"
)

const (
	defaultInboxSize   = 1024
	defaultMaxRestarts = 3
)

var defaultRestartDelay = 500 * time.Millisecond

type ReceiveFunc = func(*Context)

type MiddlewareFunc = func(ReceiveFunc) ReceiveFunc

type Opts struct {
	Producer     Producer
	Kind         string
	ID           string
	MaxRestarts  int32
	RestartDelay time.Duration
	InboxSize    int
	Middleware   []MiddlewareFunc
	Context      context.Context
}

type OptFunc func(*Opts)

// process 생성시 적용될 기본 옵션
func DefaultOpts(p Producer) Opts {
	return Opts{
		Context:      context.Background(),
		Producer:     p,
		MaxRestarts:  defaultMaxRestarts,
		InboxSize:    defaultInboxSize,
		RestartDelay: defaultRestartDelay,
		Middleware:   []MiddlewareFunc{},
	}
}

func WithContext(ctx context.Context) OptFunc {
	return func(opts *Opts) {
		opts.Context = ctx
	}
}

func WithMiddleware(mw ...MiddlewareFunc) OptFunc {
	return func(opts *Opts) {
		opts.Middleware = append(opts.Middleware, mw...)
	}
}

func WithRestartDelay(d time.Duration) OptFunc {
	return func(opts *Opts) {
		opts.RestartDelay = d
	}
}

func WithInboxSize(size int) OptFunc {
	return func(opts *Opts) {
		opts.InboxSize = size
	}
}

func WithMaxRestarts(n int) OptFunc {
	return func(opts *Opts) {
		opts.MaxRestarts = int32(n)
	}
}

func WithID(id string) OptFunc {
	return func(opts *Opts) {
		opts.ID = id
	}
}
