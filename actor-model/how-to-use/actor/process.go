package actor

import (
	"bytes"
	"fmt"
	"log/slog"
	"runtime/debug"
	"sync"
	"time"

	"github.com/DataDog/gostackparse"
)

type Envelope struct {
	Msg    any
	Sender *PID
}

// Processer is an interface the abstracts the way a process behaves.
type Processer interface {
	Start()
	PID() *PID
	Send(*PID, any, *PID)
	Invoke([]Envelope)
	Shutdown(*sync.WaitGroup)
}

type process struct {
	Opts

	inbox    Inboxer
	context  *Context
	pid      *PID
	restarts int32
	mbuffer  []Envelope
}

func newProcess(e *Engine, opts Opts) *process {
	// address와 PID를 가지고있는 PID instance 생성
	pid := NewPID(e.address, opts.Kind+pidSeparator+opts.ID)
	// 상속 구조와 메세지, 표준 컨텍스트 등을 갖는 Context 객체 생성
	ctx := newContext(opts.Context, e, pid)
	// process 인스턴스 생성
	p := &process{
		pid:     pid,
		inbox:   NewInbox(opts.InboxSize),
		Opts:    opts,
		context: ctx,
		mbuffer: nil,
	}
	p.inbox.Start(p)
	return p
}

func applyMiddleware(rcv ReceiveFunc, middleware ...MiddlewareFunc) ReceiveFunc {
	// logMiddleware(authMiddleware(ReceiveFunc))같이 ReceiveFunc에 차례대로 미들웨어 감싸기
	for i := len(middleware) - 1; i >= 0; i-- {
		rcv = middleware[i](rcv)
	}
	return rcv
}

func (p *process) Invoke(msgs []Envelope) {
	var (
		// 처리되어야할 메세지의 개수
		nmsg = len(msgs)
		// numbers of msgs that are processed.
		nproc = 0
		// FIXME: We could use nrpoc here, but for some reason placing nproc++ on the
		// bottom of the function it freezes some tests. Hence, I created a new counter
		// for bookkeeping.
		processed = 0
	)
	defer func() {
		// If we recovered, we buffer up all the messages that we could not process
		// so we can retry them on the next restart.
		// 패닉 복구 시
		if v := recover(); v != nil {
			// message를 Stopped로 변경하고
			p.context.message = Stopped{}
			// ReceiveFunc 실행
			p.context.receiver.Receive(p.context)
			// process의 메세지 버퍼에 메세지를 담고
			p.mbuffer = make([]Envelope, nmsg-nproc)
			for i := 0; i < nmsg-nproc; i++ {
				p.mbuffer[i] = msgs[i+nproc]
			}
			// process 재실행
			p.tryRestart(v)
		}
	}()
	for i := 0; i < len(msgs); i++ {
		// 처리된 메세지 개수 증가
		nproc++
		// 메세지 가져오기
		msg := msgs[i]
		// msg.Msg가 poisonPill 타입인지 검사
		// 맞다면 전체 메세지 여기서 한번에 처리(poisionPill을 제외하고 처리하기 위함)
		if pill, ok := msg.Msg.(poisonPill); ok {
			// If we need to gracefuly stop, we process all the messages
			// from the inbox, otherwise we ignore and cleanup.
			// 정상종료 설정이 true라면
			if pill.graceful {
				// 남아있는 메세지들을 차례대로 처리
				msgsToProcess := msgs[processed:]
				for _, m := range msgsToProcess {
					// ReceiveFunc 실행
					p.invokeMsg(m)
				}
			}
			// actor 삭제
			p.cleanup(pill.wg)
			return
		}
		// ReceiveFunc 실행
		p.invokeMsg(msg)
		processed++
	}
}

func (p *process) invokeMsg(msg Envelope) {
	// poisionPill 메세지는 private하여 처리하지 않음
	if _, ok := msg.Msg.(poisonPill); ok {
		return
	}
	p.context.message = msg.Msg
	p.context.sender = msg.Sender
	recv := p.context.receiver
	if len(p.Opts.Middleware) > 0 {
		// 미들웨어를 감싼 뒤 ReceiveFunc 실행
		applyMiddleware(recv.Receive, p.Opts.Middleware...)(p.context)
	} else {
		// ReceiveFunc 실행
		recv.Receive(p.context)
	}
}

func (p *process) Start() {
	recv := p.Producer()
	p.context.receiver = recv
	defer func() {
		// 패닉이 일어났을 때 비정상 종료를 방지하기 위한 부분
		if v := recover(); v != nil {
			// Stopped 메세지를 추가하고
			p.context.message = Stopped{}
			// 콜백함수 실행
			p.context.receiver.Receive(p.context)
			p.tryRestart(v)
		}
	}()
	// Initialized 메세지 추가
	p.context.message = Initialized{}
	// 미들웨어를 적용하여 ReceiveFunc 실행
	// 메세지 상태의 따라 recv.Receive 내부에서 어떤 처리를 할지 결정
	applyMiddleware(recv.Receive, p.Opts.Middleware...)(p.context)
	// 이벤트 발생
	p.context.engine.BroadcastEvent(ActorInitializedEvent{PID: p.pid, Timestamp: time.Now()})

	// Started 메시지 추가
	p.context.message = Started{}
	// 미들웨어를 적용하여 ReceiveFunc 실행
	// 메세지 상태에 따라 recv.Receive 내부에서 어떤 처리를 할지 결정
	applyMiddleware(recv.Receive, p.Opts.Middleware...)(p.context)
	// 이벤트 발생
	p.context.engine.BroadcastEvent(ActorStartedEvent{PID: p.pid, Timestamp: time.Now()})
	// 메세지 버퍼에 메세지가 있다면 메세지들의 ReceiveFunc 실행
	if len(p.mbuffer) > 0 {
		p.Invoke(p.mbuffer)
		p.mbuffer = nil
	}
}

func (p *process) tryRestart(v any) {
	// InternalError does not take the maximum restarts into account.
	// For now, InternalError is getting triggered when we are dialing
	// a remote node. By doing this, we can keep dialing until it comes
	// back up. NOTE: not sure if that is the best option. What if that
	// node never comes back up again?

	// 내부 에러일 경우 option에 적힌 delay만큼 대기한 뒤 Start() 실행
	// 내부 에러는 외부 노드에 dialing할 때 발생하기 때문에 지속적으로 dialing하기 위함
	if msg, ok := v.(*InternalError); ok {
		slog.Error(msg.From, "err", msg.Err)
		time.Sleep(p.Opts.RestartDelay)
		p.Start()
		return
	}
	stackTrace := cleanTrace(debug.Stack())
	// 최대 재시작 시도 횟수에 도달했다면 이벤트 발생 후 종료
	if p.restarts == p.MaxRestarts {
		p.context.engine.BroadcastEvent(ActorMaxRestartsExceededEvent{
			PID:       p.pid,
			Timestamp: time.Now(),
		})
		p.cleanup(nil)
		return
	}

	// 내부 에러가 아니라면 재실행 횟수를 증가시키고
	p.restarts++
	// Restart the process after its restartDelay
	// 이벤트를 발생시킨 뒤
	p.context.engine.BroadcastEvent(ActorRestartedEvent{
		PID:        p.pid,
		Timestamp:  time.Now(),
		Stacktrace: stackTrace,
		Reason:     v,
		Restarts:   p.restarts,
	})
	// 잠시 대기 후 Start() 실행
	time.Sleep(p.Opts.RestartDelay)
	p.Start()
}

func (p *process) cleanup(wg *sync.WaitGroup) {
	// 부모 컨텍스트가 있으면 자신을 자식 목록에서 제거
	if p.context.parentCtx != nil {
		p.context.parentCtx.children.Delete(p.Kind)
	}

	// 자식 액터가 있으면 각각에게 중지 신호를 보내기
	if p.context.children.Len() > 0 {
		children := p.context.Children()
		for _, pid := range children {
			p.context.engine.Poison(pid).Wait()
		}
	}

	// 상태를 stopped로 변경
	p.inbox.Stop()
	p.context.engine.Registry.Remove(p.pid)
	p.context.message = Stopped{}
	applyMiddleware(p.context.receiver.Receive, p.Opts.Middleware...)(p.context)

	// actor가 종료되었다는 이벤트 발생
	p.context.engine.BroadcastEvent(ActorStoppedEvent{PID: p.pid, Timestamp: time.Now()})
	if wg != nil {
		wg.Done()
	}
}

func (p *process) PID() *PID { return p.pid }
func (p *process) Send(_ *PID, msg any, sender *PID) {
	p.inbox.Send(Envelope{Msg: msg, Sender: sender})
}
func (p *process) Shutdown(wg *sync.WaitGroup) { p.cleanup(wg) }

func cleanTrace(stack []byte) []byte {
	goros, err := gostackparse.Parse(bytes.NewReader(stack))
	if err != nil {
		slog.Error("failed to parse stacktrace", "err", err)
		return stack
	}
	if len(goros) != 1 {
		slog.Error("expected only one goroutine", "goroutines", len(goros))
		return stack
	}
	// skip the first frames:
	goros[0].Stack = goros[0].Stack[4:]
	buf := bytes.NewBuffer(nil)
	_, _ = fmt.Fprintf(buf, "goroutine %d [%s]\n", goros[0].ID, goros[0].State)
	for _, frame := range goros[0].Stack {
		_, _ = fmt.Fprintf(buf, "%s\n", frame.Func)
		_, _ = fmt.Fprint(buf, "\t", frame.File, ":", frame.Line, "\n")
	}
	return buf.Bytes()
}
