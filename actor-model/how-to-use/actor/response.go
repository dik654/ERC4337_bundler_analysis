package actor

import (
	"context"
	"math"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type Response struct {
	engine  *Engine
	pid     *PID
	result  chan any
	timeout time.Duration
}

func NewResponse(e *Engine, timeout time.Duration) *Response {
	return &Response{
		engine:  e,
		result:  make(chan any, 1),
		timeout: timeout,
		pid:     NewPID(e.address, "response"+pidSeparator+strconv.Itoa(rand.Intn(math.MaxInt32))),
	}
}

func (r *Response) Result() (any, error) {
	// context 제한시간 설정
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	// 종료시 registry에서 actor 제거
	defer func() {
		cancel()
		r.engine.Registry.Remove(r.pid)
	}()

	// 채널로 result를 받거나 제한시간 종료 신호를 받았다면 결과 리턴
	select {
	case resp := <-r.result:
		return resp, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (r *Response) Send(_ *PID, msg any, _ *PID) {
	r.result <- msg
}

func (r *Response) PID() *PID                  { return r.pid }
func (r *Response) Shutdown(_ *sync.WaitGroup) {}
func (r *Response) Start()                     {}
func (r *Response) Invoke([]Envelope)          {}
