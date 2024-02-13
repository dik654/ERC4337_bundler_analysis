package actor

import (
	"context"
	"log/slog"
)

// eventSub is the message that will be send to subscribe to the event stream.
type eventSub struct {
	pid *PID
}

// EventUnSub is the message that will be send to unsubscribe from the event stream.
type eventUnsub struct {
	pid *PID
}

type eventStream struct {
	subs map[*PID]bool
}

func newEventStream() Producer {
	return func() Receiver {
		return &eventStream{
			subs: make(map[*PID]bool),
		}
	}
}

// Receive for the event stream. All system-wide events are sent here.
// Some events are specially handled, such as eventSub, EventUnSub (for subscribing to events),
// DeadletterSub, DeadletterUnSub, for subscribing to DeadLetterEvent
func (e *eventStream) Receive(c *Context) {
	switch msg := c.Message().(type) {
	// msg type이 eventSub라면
	case eventSub:
		// subs에 추가
		e.subs[msg.pid] = true
	// 구독 해제라면
	case eventUnsub:
		// subs에서 삭제
		delete(e.subs, msg.pid)
	default:
		// check if we should log the event, if so, log it with the relevant level, message and attributes
		// EventLogger타입의 메세지라면
		logMsg, ok := c.Message().(EventLogger)
		if ok {
			// 메세지 로깅
			level, msg, attr := logMsg.Log()
			slog.Log(context.Background(), level, msg, attr...)
		}
		for sub := range e.subs {
			c.Forward(sub)
		}
	}
}
